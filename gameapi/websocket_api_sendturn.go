package gameapi

import (
    "encoding/json"

    ws "GoChessgameServer/websocket"
    u "GoChessgameServer/util"
    "GoChessgameServer/auth"
    "GoChessgameServer/game"

    "github.com/gorilla/websocket"
)

func WebsocketAPISendTurnHandle(wc *ws.WebsocketConnection, data *WebsocketRequest, tokenData *auth.JWTUserClaim) {

    conn := wc.GetConnection()
    request := game.Turn{}

    if err := json.Unmarshal(data.Args, &request); err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Invalid arg object passed"))
        return
    }

    // Get game session by game id, them get appropriate client connection
    // and deliver request
    gameSession, err := game.SessionStore.GetGameSession(request.GameID)
    if err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Game session not found by given game id"))
        return
    }

    for _, v := range gameSession.GetAllPlayers() {
        if v.Connection == wc {
            v.ReadChannel <- &request
            return
        }
    }

    conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Appropriate client connection not found in game session"))
}
