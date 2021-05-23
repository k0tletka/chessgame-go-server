package gameapi

import (
    "net/http"
    "time"

    "github.com/gorilla/websocket"
)

var (
    // Websocket upgader, used to ugrade HTTP request to long-life WS connections
    wUpgrader = websocket.Upgrader{
        HandshakeTimeout: 5 * time.Second,
        CheckOrigin: func (r *http.Request) bool { return true; },
    }
)

// This function performs websocket upgrading and connection handling
func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
    c, err := wUpgrader.Upgrade(w, r, nil)

    if err != nil {
        gameApiLogger.Printf("Error when trying to upgrade WS connection: %s\n", err.Error())
        return
    }

    // Create new WebsocketConnection and insert into list
    wc := NewWebsocketConnection(c, websocketReadHandler)
    wsStore.InsertConnection(wc)
}

// Handler for websocket requests from peer
func websocketReadHandler(wc *WebsocketConnection, data []byte) {
    // Send hello response
    response := "Input string: " + string(data)
    wc.GetConnection().WriteMessage(websocket.TextMessage, []byte(response))
}
