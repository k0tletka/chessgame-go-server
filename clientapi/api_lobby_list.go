package clientapi

import (
    "net/http"
    "encoding/json"

    u "GoChessgameServer/util"
    "GoChessgameServer/game"
)

func LobbyList(w http.ResponseWriter, r *http.Request) {

    writeError := u.WriteErrorCreator(w)
    contextUser := r.Context().Value("login").(string)

    // Prepare for response
    type responseType struct {
        GameID          int         `json:"id"`
        GameTitle       string      `json:"title"`
        Players         []string    `json:"players"`
    }

    response := []responseType{}

    for _, v := range game.SessionStore.GetAllGameSessions() {
        players := []string{}

        for _, p := range v.GetAllPlayers() {
            players = append(players, p.Login)
        }

        response = append(response, responseType{
            GameID: v.GameID,
            GameTitle: v.GameTitle,
            Players: players,
        })
    }

    if err := json.NewEncoder(w).Encode(&response); err != nil {
        writeError("Server error")
        w.WriteHeader(http.StatusInternalServerError)
        clientApiLogger.Printf("LobbyList: Error when sending response: %s\n", err.Error())
        return
    }
    w.Header().Add("Content-Type", "application/json")

    // Log
    clientApiLogger.Printf("LobbyList: User %s requested lobby list\n", contextUser)
}
