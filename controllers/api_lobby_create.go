package controllers

import (
    "net/http"
    "encoding/json"
    "reflect"

    u "GoChessgameServer/util"
    "GoChessgameServer/store"
)

// This function creates new game and
// pushes it inside lobby store
func LobbyCreate(w http.ResponseWriter, r *http.Request) {

    contextUser := r.Context().Value("login").(string)
    writeError := u.WriteErrorCreator(w)

    type reqType struct{
        GameTitle string `json:"title"`
    }
    req := reqType{}

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || reflect.DeepEqual(req, reqType{}) {
        writeError("Invalid request")
        contrLogger.Printf("LobbyCreate: Error when parsing request from client: %s\n", err.Error())
        return
    }

    // Create new game and insert it
    gameId := store.RegisterNewGameStore(req.GameTitle, contextUser)

    // Response with game id
    resp := struct{
        GameID int `json:"id"`
    }{
        GameID: gameId,
    }

    // Send response
    if err := json.NewEncoder(w).Encode(resp); err != nil {
        writeError("Server error")
        w.WriteHeader(http.StatusInternalServerError)
        contrLogger.Printf("LobbyCreate: Error when sending response: %s\n", err.Error())
        return
    }
    w.Header().Add("Content-Type", "application/json")

    // Log
    contrLogger.Printf("LobbyCreate: User %s created a new game with id %d\n", contextUser, gameId)
}
