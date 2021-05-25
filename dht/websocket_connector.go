package dht

import (
    "encoding/json"
    "net/http"
    "time"

    ws "GoChessgameServer/websocket"
    u "GoChessgameServer/util"

    "github.com/gorilla/websocket"
)

// This type represents base request
type dhtAPIBaseRequest struct {
    MethodName  string          `json:"method_name"`
    Args        json.RawMessage `json:"args"`
}

var (
    // Websocket upgrader object
    wUpgrader = &websocket.Upgrader{
        HandshakeTimeout: 5 * time.Second,
        CheckOrigin: func (r *http.Request) bool { return true; },
    }
)

// This functions performs websocket upgrading and connection handling
func (m *DHTManager) websocketHandler(w http.ResponseWriter, r *http.Request) {
    c, err := wUpgrader.Upgrade(w, r, nil)

    if err != nil {
        dhtLogger.Printf("Error when trying to upgrade WS connection: %s\n", err.Error())
        return
    }

    // Create new websocket connection
    ws.NewWebsocketConnection(c, m.connectionServerReadHandler, m.WsServerConns)
}

// Main handler for websocker requests
func (m *DHTManager) connectionServerReadHandler(wc *ws.WebsocketConnection, data []byte) {
    conn := wc.GetConnection()
    request := dhtAPIBaseRequest{}

    if err := json.Unmarshal(data, &request); err != nil {
        // Try to parse into error object
        errorObject := struct{
            Error   string  `json:"error"`
        }{}

        if err := json.Unmarshal(data, &errorObject); err != nil {
            dhtLogger.Printf("Error when trying to unmarshal response object from %s: %s\n", conn.RemoteAddr(), err.Error())
            return
        }

        dhtLogger.Printf("Error response from %s: %s\n", conn.RemoteAddr(), errorObject.Error)
        return
    }

    // Route requests
    routingPaths := map[string]func(*ws.WebsocketConnection, []byte){
        "handshake": m.handshakeMethodServerHandler,
    }

    if handler, ok := routingPaths[request.MethodName]; ok {
        handler(wc, request.Args)
    } else {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("There is no handler for provided method"))
    }
}
