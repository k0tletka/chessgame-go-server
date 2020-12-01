package controllers

import (
    "net/http"
    "encoding/json"

    u "GoChessgameServer/util"
    "GoChessgameServer/store"
)

// This function returns a list of current
// waiting games
func LobbyList(w http.ResponseWriter, r *http.Request) {

    contextUser := r.Context().Value("login").(string)
    writeError := u.WriteErrorCreator(w)

    type respType struct{
        GameID int `json:"id"`
        PlayerName string `json:"player"`
        GameTitle string `json:"title"`
    }
    resp := []respType{}

    for _, game := range store.GameStores {
        resp = append(resp, respType{
            GameID: game.GameID,
            PlayerName: game.PlayerOneLogin,
            GameTitle: game.GameTitle,
        })
    }

    // Send response
    if err := json.NewEncoder(w).Encode(resp); err != nil {
        writeError("Server error")
        w.WriteHeader(http.StatusInternalServerError)
        contrLogger.Printf("LobbyList: Error when sending response: %s\n", err.Error())
        return
    }
    w.Header().Add("Content-Type", "application/json")

    // Log new user
    contrLogger.Printf("LobbyList: User %s requested lobby list\n", contextUser)
}
