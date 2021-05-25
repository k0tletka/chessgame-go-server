package dht

import (
    "net/http"
    "fmt"
    "context"
    "time"
    "encoding/json"
    "encoding/hex"

    ws "GoChessgameServer/websocket"
    c "GoChessgameServer/conf"
    u "GoChessgameServer/util"

    "github.com/gorilla/websocket"
)

// This type represent static peer from configuration
type StaticPeerConnection struct {
    Connection *ws.WebsocketConnection
    StaticPeer *c.SPeer
}

// This is main manager of server DHT API,
// that performs handshaking process and provides
// programm API for other modules
type DHTManager struct {
    // Store of server connections
    WsServerConns *ws.WebsocketStore

    // Store of client connections
    WsClientConns *ws.WebsocketStore

    // List with connect information about static peers and connection object
    staticPeerConnections []*StaticPeerConnection

    // Websocket dialer
    wsDialer *websocket.Dialer
}

func CreateNewDHTManager() *DHTManager {
    wsDialer := &websocket.Dialer{
        HandshakeTimeout: 5 * time.Second,
    }

    newManager := &DHTManager {
        WsServerConns: &ws.WebsocketStore{},
        WsClientConns: &ws.WebsocketStore{},
        staticPeerConnections: []*StaticPeerConnection{},
        wsDialer: wsDialer,
    }

    // Fill static peer connections
    newManager.fillStaticPeerConnections()

    // Start handshake goroutine
    go newManager.startHandshakeProcedure()

    return newManager
}

func (m *DHTManager) fillStaticPeerConnections() {
    for _, v := range c.Conf.StaticPeers {
        m.staticPeerConnections = append(m.staticPeerConnections, &StaticPeerConnection{
            Connection: nil,
            StaticPeer: &v,
        })
    }
}

func (m *DHTManager) startHandshakeProcedure() {
    var handshakeTimeout uint

    if !c.DecodeMetadata.IsDefined("dht_api", "handshake_period") {
        handshakeTimeout = 300 // 5 minutes default
    } else {
        handshakeTimeout = c.Conf.DHTApi.HandshakePeriod
    }

    for {
        // Enstablish connection with static peers and send them request with out
        // server identifier
        for _, v := range m.staticPeerConnections {
            var conn *ws.WebsocketConnection

            // If connection is nil or connection closed, create new connection
            if v.Connection == nil || v.Connection.Closed() {
                wsConn, err := m.createConnection(v)

                if err != nil {
                    dhtLogger.Printf("Connection to %s peer failed: %s\n", v.StaticPeer.ServerName, err.Error())
                    continue
                }

                // Add connection to list
                conn = ws.NewWebsocketConnection(wsConn, m.connectionClientReadHandler, m.WsClientConns)
                v.Connection = conn
            } else {
                conn = v.Connection
            }

            // Get listening port of server api
            _, listenport := getListenInformation()

            // Connection exist, so just send handshake message
            request := struct{
                ServerIdentifier    string  `json:"server_identifier"`
                ServerAPIPort       uint16  `json:"server_api_port"`
            }{
                ServerIdentifier: hex.EncodeToString(dhtServerIdentifier[:]),
                ServerAPIPort: listenport,
            }

            var data []byte
            var err error

            if data, err = json.Marshal(&request); err != nil {
                dhtLogger.Printf("Error when marshalling request to peer %s: %s\n", v.StaticPeer.ServerName, err.Error())
                continue
            }

            baseRequest := dhtAPIBaseRequest{
                MethodName: "handshake",
                Args: data,
            }

            if data, err = json.Marshal(&baseRequest); err != nil {
                dhtLogger.Printf("Error when marshalling request to peer %s: %s\n", v.StaticPeer.ServerName, err.Error())
                continue
            }

            conn.GetConnection().WriteMessage(websocket.TextMessage, data)
        }

        // Sleep until next handshake timeout
        time.Sleep(time.Duration(handshakeTimeout) * time.Second)
    }
}

func (m *DHTManager) createConnection(s *StaticPeerConnection) (*websocket.Conn, error) {
    // Connect with timeout
    var timeout uint

    if !c.DecodeMetadata.IsDefined("dht_api", "peer_connection_timeout") {
        timeout = 5 // 5 seconds by default
    } else {
        timeout = c.Conf.DHTApi.PeerConnTimeout
    }

    context, cancel := context.WithTimeout(context.Background(), time.Duration(timeout) * time.Second)
    defer cancel()

    conn, _, err := m.wsDialer.DialContext(context, fmt.Sprintf(
        "ws%s://%s:%d/ws",
        map[bool]string{true: "s", false: ""}[s.StaticPeer.UseTLS],
        s.StaticPeer.ServerName,
        s.StaticPeer.ConnectionPort,
    ), http.Header{})

    return conn, err
}

// Main read handler for connections
func (m *DHTManager) connectionClientReadHandler(wc *ws.WebsocketConnection, data []byte) {
    conn := wc.GetConnection()
    response := dhtAPIBaseRequest{}

    if err := json.Unmarshal(data, &response); err != nil {
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
        "handshake": m.handshakeMethodClientHandler,
    }

    if handler, ok := routingPaths[response.MethodName]; ok {
        handler(wc, response.Args)
    } else {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("There is no handler for provided method"))
    }
}

// Client handler of handshake request
func (m *DHTManager) handshakeMethodClientHandler(wc *ws.WebsocketConnection, data []byte) {
}
