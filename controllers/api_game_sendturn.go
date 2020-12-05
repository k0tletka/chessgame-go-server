package controllers

import (
    "net/http"
    "encoding/json"
    "reflect"

    u "GoChessgameServer/util"
    "GoChessgameServer/store"
)

// This function allows to send player turn
// to game thread and retrieve operation results
func SendTurn(w http.ResponseWriter, r *http.Request) {

    contextUser := r.Context().Value("login").(string)
    writeError := u.WriteErrorCreator(w)

    type reqType struct {
        GameID int `json:"id"`
        Surrender bool `json:"surrender"`
        Figposx int `json:"figposx"`
        Figposy int `json:"figposy"`
        AltX int `json:"xalt"`
        AltY int `json:"yalt"`
    }
    req := reqType{}

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || reflect.DeepEqual(req, reqType{}) {
        writeError("Invalid request")
        contrLogger.Printf("SendTurn: Error when parsing request from client: %s\n", err.Error())
        return
    }

    // Get game store
    gameStore, err := store.GetGameStore(req.GameID)
    if err != nil {
        writeError("Game with the specified id is not found")
        return
    }

    // Check, that user must send turn
    if (gameStore.IsPlayerOneTurn && contextUser != gameStore.PlayerOneLogin) || (!gameStore.IsPlayerOneTurn && contextUser != gameStore.PlayerTwoLogin) {
        writeError("You cant send turn at this state")
        contrLogger.Printf("SendTurn: Player %s tried to send turn when he is not allowed to do so\n", contextUser)
        return
    }

    // Send turn to game thread and read response
    gameTurn := store.GameTurn{
        FigposX: req.Figposx,
        FigposY: req.Figposy,
        AltX: req.AltX,
        AltY: req.AltY,
        Surrender: req.Surrender,
    }
    gameStore.SendTurnRequest <- &gameTurn
    response := <-gameStore.SendTurnResponse

    if (response != "") {
        writeError(response)
        return
    }

    // All ok
    contrLogger.Printf("SendTurn: User %s sended his turn on game id %d\n", contextUser, req.GameID)
}
