package dht

import (
    "time"
    "encoding/json"
    "encoding/hex"
    "net"

    ws "GoChessgameServer/websocket"
    u "GoChessgameServer/util"
    "GoChessgameServer/database"

    "github.com/gorilla/websocket"
)

func (m *DHTManager) handshakeMethodHandler(wc *ws.WebsocketConnection, data *dhtAPIBaseRequest) {
    conn := wc.GetConnection()

    // Get server request identifier
    request := struct{
        ServerIdentifier    string  `json:"server_identifier"`
        ServerAPIPort       uint16  `json:"server_api_port"`
        UseTLS              bool    `json:"server_api_use_tls"`
        ConnectionLimit     *uint   `json:"connection_limit,omitempty"`
        ConnectingStatic    bool    `json:"connecting_static"`
        IsPeerStatic        bool    `json:"is_peer_static"`
    }{}

    if err := json.Unmarshal(data.Args, &request); err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Invalid request"))
        return
    }

    // Memory information about host, who maked request
    decodedIdentifier, err := hex.DecodeString(request.ServerIdentifier)
    if err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Invalid server identifier has passed"))
        return
    }

    hostInfo := database.DHTHosts{
        ServerIdentifier: decodedIdentifier,
        SrvLocalIdentifier: dhtServerIdentifier[:],
        IPAddress: conn.RemoteAddr().(*net.TCPAddr).IP.String(),
        Port: request.ServerAPIPort,
        UseTLS: request.UseTLS,
        IsPeerStatic: request.IsPeerStatic,
        IsPeerConnsStatic: request.ConnectingStatic,
        LastHandshake: time.Now(),
    }

    database.DB.Save(&hostInfo)

    // Add to handshake map to send handshake requests later
    if _, ok := m.databasePeerConnections[request.ServerIdentifier]; !ok {
        m.databasePeerConnections[request.ServerIdentifier] = wc
    }

    // Results to requester
    type resType struct {
        ServerIdentifier    string  `json:"server_identifier"`
        ServerAPIPort       uint16  `json:"server_api_port"`
        ServerIPAddress     string  `json:"server_api_ip_address"`
        UseTLS              bool    `json:"server_api_use_tls"`
    }

    results := []resType{}

    // Prepare hosts list, that are places near by hash
    serverList := []database.DHTHosts{}

    // Select all dht hosts
    database.DB.Find(&serverList)

    for _, v := range serverList {
        // Currenly, all instances stores maximum public instances in their databases
        /*if request.ConnectionLimit != nil && int(*(request.ConnectionLimit)) >= len(results) {
            break
        }*/

        if hex.EncodeToString(v.ServerIdentifier) == request.ServerIdentifier {
            continue
        }

        results = append(results, resType{
            ServerIdentifier: hex.EncodeToString(v.ServerIdentifier),
            ServerAPIPort: v.Port,
            ServerIPAddress: v.IPAddress,
            UseTLS: v.UseTLS,
        })
    }

    var resData []byte

    if resData, err = json.Marshal(&results); err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Error occured when marshalling request: " + err.Error()))
        return
    }

    // Send result back to client
    response := dhtAPIBaseRequest{
        MethodName: "handshake_response",
        Args: resData,
    }

    if resData, err = json.Marshal(&response); err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Error occured when marshalling request: " + err.Error()))
        return
    }

    conn.WriteMessage(websocket.TextMessage, resData)
}
