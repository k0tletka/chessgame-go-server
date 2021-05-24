package gameapi

import (
    "encoding/json"
    "fmt"

    "GoChessgameServer/auth"
    ws "GoChessgameServer/websocket"
    u "GoChessgameServer/util"
    "GoChessgameServer/game"

    "github.com/gorilla/websocket"
)

// Handle game connections
func WebsocketAPIConnectHandle(wc *ws.WebsocketConnection, data *WebsocketRequest, tokenData *auth.JWTUserClaim) {

    conn := wc.GetConnection()

    request := struct{
        GameID   int  `json:"game_id"`
    }{}

    if err := json.Unmarshal(data.Args, &request); err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Invalid arg object passed"))
        return
    }

    // Connection to existing game session
    gameSession, err := game.SessionStore.GetGameSession(request.GameID)
    if err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Session with given game id not exist"))
        return
    }

    gameClientConnection := &game.GameClientConnection{
        Connection: wc,
        Login: tokenData.Login,
        ExternalUser: data.FromExternalInstance,
        ReadChannel: make(chan *game.Turn, 1),
    }

    err = gameSession.AddNewConnection(gameClientConnection)
    if err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson(fmt.Sprintf("Error when connecting to game: %s\n", err.Error())))
        return
    }
}
