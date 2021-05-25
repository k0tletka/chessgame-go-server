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

func (m *DHTManager) handshakeMethodServerHandler(wc *ws.WebsocketConnection, data []byte) {
    conn := wc.GetConnection()

    // Get server request identifier
    request := struct{
        ServerIdentifier    string  `json:"server_identifier"`
        ServerAPIPort       uint16  `json:"server_api_port"`
        UseTLS              bool    `json:"server_api_use_tls"`
        ConnectionLimit     *uint   `json:"connection_limit,omitempty"`
    }{}

    if err := json.Unmarshal(data, &request); err != nil {
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
        LastHandshake: time.Now(),
    }

    database.DB.Save(&hostInfo)

    // Results to requester
    type resType struct {
        ServerIdentifier    string  `json:"server_identifier"`
        ServerAPIPort       uint16  `json:"server_api_port"`
        UseTLS              bool    `json:"server_api_use_tls"`
    }

    results := []resType{}

    // Prepare hosts list, that are places near by hash
    serverList := []database.DHTHosts{}

    // Select all dht hosts
    database.DB.Find(&serverList)

    for _, v := range serverList {
        if request.ConnectionLimit != nil && int(*(request.ConnectionLimit)) >= len(results) {
            break
        }

        results = append(results, resType{
            ServerIdentifier: hex.EncodeToString(v.ServerIdentifier),
            ServerAPIPort: v.Port,
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
        MethodName: "handshake",
        Args: resData,
    }

    if resData, err = json.Marshal(&response); err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Error occured when marshalling request: " + err.Error()))
        return
    }

    conn.WriteMessage(websocket.TextMessage, data)
}
