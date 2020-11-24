package store

import (
    "crypto/rand"
    "log"
    "os"

    "GoChessgameServer/logger"

    jwt "github.com/dgrijalva/jwt-go"
)

// Store logger
var storeLogger *log.Logger

// This string contains a key, stat will be
// used for jwt signing
var JWTKey []byte

// This type represent token claims with login string
type JWTClaims struct {
    Login string
    jwt.StandardClaims
}

func init() {
    // Create store logger
    storeLogger = logger.AddNewLogger("Store", os.Stdout, log.LstdFlags | log.Lmsgprefix)

    // Generate JWT key
    JWTKey = make([]byte, 256, 256)
    _, err := rand.Read(JWTKey)
    if err != nil {
        storeLogger.Fatalf("Error when generating key: %s", err.Error())
    }
}
