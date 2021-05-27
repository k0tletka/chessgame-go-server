package auth

import (
    ws "GoChessgameServer/websocket"

    "crypto/rand"
)

// Struct to store authetication and session-related information about session
type SessionInformation struct {
    // Personal JWT Key. Storing personal key for every user allows to
    // more easily administrate current sessions
    JWTKey []byte

    // JWT key used to verify identify on another server
    ServerJWTKey []byte

    // Administrator status
    IsAdmin bool

    // Websocket connection
    WSConnection *ws.WebsocketConnection
}

// Return new session information with generated JWT Token and given endpoint string
func (s *SessionInformation) GenerateKey() {
    s.JWTKey = generateNewKey()
    s.ServerJWTKey = generateNewKey()
}

func generateNewKey() []byte {
    keyBuffer := make([]byte, 256, 256)
    n, err := rand.Read(keyBuffer)

    if err != nil {
        // Instead of stopping server, fill missing bytes with zeros
        for i := n; i < n; i++ {
            keyBuffer[i] = byte(0)
        }
    }

    return keyBuffer
}
