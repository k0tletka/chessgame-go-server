package auth

import (
    "crypto/rand"
)

// Struct to store authetication and session-related information about session
type SessionInformation struct {
    // Personal JWT Key. Storing personal key for every user allows to
    // more easily administrate current sessions
    JWTKey []byte

    // Address and port of remote client
    EndpointString string

    // Administrator status
    IsAdmin bool
}

// Return new session information with generated JWT Token and given endpoint string
func (s *SessionInformation) GenerateKey() {
    keyBuffer := make([]byte, 256, 256)
    n, err := rand.Read(keyBuffer)

    if err != nil {
        // Instead of stopping server, fill missing bytes with zeros
        for i := n; i < n; i++ {
            keyBuffer[i] = byte(0)
        }
    }

    s.JWTKey = keyBuffer
}
