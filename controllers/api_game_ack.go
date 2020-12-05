package controllers

import (
    "net/http"
    "encoding/json"
    "reflect"

    u "GoChessgameServer/util"
    "GoChessgameServer/store"
)

// This function allows users
// to acknowledge themselves,
// that they are on the line
func GameAck(w http.ResponseWriter, r *http.Request) {

    contextUser := r.Context().Value("login").(string)
    writeError := u.WriteErrorCreator(w)

    type reqType struct{
        GameID int `json:"id"`
    }
    req := reqType{}

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || reflect.DeepEqual(req, reqType{}) {
        writeError("Invalid request")
        contrLogger.Printf("GameAck: Error when parsing request from client: %s\n", err.Error())
        return
    }

    // Get game store
    gameStore, err := store.GetGameStore(req.GameID)
    if err != nil {
        return
    }

    // Check that user can ack on this game
    if contextUser != (*gameStore).PlayerOneLogin && contextUser != (*gameStore).PlayerTwoLogin {
        writeError("You can't acknowledge on this game")
        return
    }

    // Make acknowledge
    (*gameStore).AckChannel <- contextUser

    // Log acknowledge
    contrLogger.Printf("GameAck: User %s is acknowledged longpoll response for game id %d\n", contextUser, req.GameID)
}
