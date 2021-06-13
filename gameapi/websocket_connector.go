package gameapi

import (
    "net/http"
    "time"
    "encoding/json"

    "GoChessgameServer/auth"
    ws "GoChessgameServer/websocket"
    u "GoChessgameServer/util"
    "GoChessgameServer/dht"

    "github.com/gorilla/websocket"
)

// Base type for websocket data
type WebsocketRequest struct {
    Token                   string          `json:"token"`
    MethodPath              string          `json:"method_path"`

    // This field means, that message has came from client of another instance,
    // and we must verify its token by sending request to appropriate instance.
    FromExternalInstance    bool            `json:"from_external_instance"`
    ServerInstanceIdentify  *string         `json:"server_instance_identify,omitempty"`

    // Arguments that passed to method handler function.
    // Content depends on method handler's needs.
    Args                    json.RawMessage `json:"args"`
}


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

    // Create new WebsocketConnection
    ws.NewWebsocketConnection(c, websocketReadHandler, wsStore)
}

// Handler for websocket requests from peer
func websocketReadHandler(wc *ws.WebsocketConnection, data []byte) {
    // Connection for result writing
    conn := wc.GetConnection()

    // Parse request
    req := WebsocketRequest{}

    if err := json.Unmarshal(data, &req); err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Invalid request"))
        return
    }

    var tokenData *auth.JWTUserClaim

    if req.FromExternalInstance {
        if req.ServerInstanceIdentify == nil {
            conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Server instance identity can't be empty"))
            return
        }

        // Send request to verify server token
        var verified bool
        var err error

        tokenData, verified, err = dht.DHTMgr.VerifyServerToken(req.Token, *req.ServerInstanceIdentify)

        if err != nil {
            conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Can't verify server token: " + err.Error()))
            return
        }

        if !verified {
            conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Invalid token"))
            return
        }
    } else {
        // Verify jwt token
        var verified bool
        tokenData, verified = auth.VerifyToken(req.Token)

        if !verified {
            conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Invalid token"))
            gameApiLogger.Printf("Invalid token passed from %s\n", conn.RemoteAddr())
            return
        }
    }

    // Perform request routing based on request method
    routingPaths := map[string]func(*ws.WebsocketConnection, *WebsocketRequest, *auth.JWTUserClaim){
        "api_connect": WebsocketAPIConnectHandle,
        "api_create": WebsocketAPICreateHandle,
        "api_sendturn": WebsocketAPISendTurnHandle,
    }

    if handler, ok := routingPaths[req.MethodPath]; ok {
        handler(wc, &req, tokenData)
    } else {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("There is no handler for provided method"))
        return
    }
}
