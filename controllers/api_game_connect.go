package controllers

import (
    "net/http"
    "encoding/json"
    "time"
    "reflect"
    "strconv"
    "sync"

    u "GoChessgameServer/util"
    "GoChessgameServer/store"
    "GoChessgameServer/game"
)

// This mutex allow to sync connect, so only the first player,
// who sent request, will be connected to the game
var connectSync = sync.Mutex{}

// This function performs connection
// to the existent game and start game session
func GameConnect(w http.ResponseWriter, r *http.Request) {

    contextUser := r.Context().Value("login").(string)
    writeError := u.WriteErrorCreator(w)

    type reqType struct {
        GameID int `json:"id"`
    }
    req := reqType{}

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || reflect.DeepEqual(req, reqType{}) {
        writeError("Invalid request")
        contrLogger.Printf("GameConnect: Error when parsing request from client: %s\n", err.Error())
        return
    }

    // Get game store
    gameStore, err := store.GetGameStore(req.GameID)
    if err != nil {
        writeError("Game with the specified id is not found")
        return
    }

    connectSync.Lock()
    // Check self connect and double connect
    if (*gameStore).PlayerOneLogin == contextUser {
        writeError("You can't connect to yourself")
        return
    } else if (*gameStore).PlayerTwoLogin != "" {
        writeError("Some player connected already")
        return
    }

    // Set second player and game started flag
    (*gameStore).PlayerTwoLogin = contextUser
    (*gameStore).GameStarted = true
    connectSync.Unlock()

    resplm := struct {
        Login string `json:"opponent"`
    }{
        Login: contextUser,
    }
    jsonBytes, err := json.Marshal(&resplm)
    if err != nil {
        writeError("Server error")
        w.WriteHeader(http.StatusInternalServerError)
        contrLogger.Printf("GameConnect: Error when marshaling json: %s\n", err.Error())
        return
    }

    // Notify first player that opponent finded
    store.WaitGameLM.Publish(strconv.Itoa(req.GameID), string(jsonBytes))

    // Wait for 3 seconds
    timeoutTimer := time.NewTimer(time.Second * 3)

    select {
    case <-timeoutTimer.C:
        writeError("Connection timed out, aborting")
        _ = store.RemoveGameStore(req.GameID)
        return
    case login := <-gameStore.AckChannel:
        if login == (*gameStore).PlayerOneLogin {
            resp := struct{
                Login string `json:"opponent"`
            }{
                Login: login,
            }

            // Send response
            if err := json.NewEncoder(w).Encode(resp); err != nil {
                writeError("Server error")
                w.WriteHeader(http.StatusInternalServerError)
                contrLogger.Printf("GameConnect: Error when sending response: %s\n", err.Error())
                return
            }

            w.Header().Add("Content-Type", "application/json")
            contrLogger.Printf("GameConnect: User %s connected to a game %d\n", contextUser, req.GameID)

            // Game process approved, start game controller in other thread
            go game.ControlGame(gameStore)
        } else {
            writeError("Game aborted, other user tried to connect to the game")
            _ = store.RemoveGameStore(req.GameID)
        }
    }
}
