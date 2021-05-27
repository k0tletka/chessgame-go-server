package dht

import (
    "encoding/json"
    "encoding/hex"
    "time"

    ws "GoChessgameServer/websocket"
    u "GoChessgameServer/util"
    "GoChessgameServer/database"

    "github.com/gorilla/websocket"
)

func (m *DHTManager) handshakeResponseMethodHandler(wc *ws.WebsocketConnection, data *dhtAPIBaseRequest) {
    conn := wc.GetConnection()

    request := []struct{
        ServerIdentifier    string  `json:"server_identifier"`
        ServerAPIPort       uint16  `json:"server_api_port"`
        ServerIPAddress     string  `json:"server_api_ip_address"`
        UseTLS              bool    `json:"use_tls"`
    }{}

    if err := json.Unmarshal(data.Args, &request); err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Invalid request"))
        return
    }

    // Save given results to database
    for _, v := range request {
        srvIdentifier, err := hex.DecodeString(v.ServerIdentifier)

        if err != nil {
            dhtLogger.Printf("Warning: server identifier for %s host failed to decode\n", conn.RemoteAddr())
            continue
        }

        databaseRecord := database.DHTHosts{
            ServerIdentifier: srvIdentifier,
            SrvLocalIdentifier: dhtServerIdentifier[:],
            IPAddress: v.ServerIPAddress,
            Port: v.ServerAPIPort,
            UseTLS: v.UseTLS,
            LastHandshake: time.Now(),
        }

        database.DB.Save(&databaseRecord)
    }
}
