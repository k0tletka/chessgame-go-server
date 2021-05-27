package gameapi

import (
    "encoding/json"

    ws "GoChessgameServer/websocket"
    u "GoChessgameServer/util"
    "GoChessgameServer/auth"
    "GoChessgameServer/game"

    "github.com/gorilla/websocket"
)

func WebsocketAPICreateHandle(wc *ws.WebsocketConnection, data *WebsocketRequest, tokenData *auth.JWTUserClaim) {

    conn := wc.GetConnection()

    request := struct{
        GameTitle   string  `json:"title"`
    }{}

    if err := json.Unmarshal(data.Args, &request); err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Invalid arg object passed"))
        return
    }

    // Creating new game session
    gameClientConnection := game.GameClientConnection{
        Connection: wc,
        Login: tokenData.Login,
        ExternalUser: data.FromExternalInstance,
        ReadChannel: make(chan *game.Turn, 1),
    }

    gameId, err := game.SessionStore.RegisterNewGameSession(
        request.GameTitle,
        &gameClientConnection,
        2,
    )

    if err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Invalid max player variable was provided"))
        return
    }

    response := struct{
        GameID  int `json:"game_id"`
    }{
        GameID: gameId,
    }

    if respData, err := json.Marshal(&response); err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Can't write response to client: " + err.Error()))
        return
    } else {
        conn.WriteMessage(websocket.TextMessage, respData)
    }
}
