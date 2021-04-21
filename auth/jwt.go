package auth

import (
    "strings"

    jwt "github.com/dgrijalva/jwt-go"
)

var (
    // SignindMethod, that used for signing and verifying tokens
    signMethod = jwt.GetSigningMethod("HS256")

    // Utility
    tokenParser = &jwt.Parser{}
)

// Type of JWT token claim, that containing user login
type JWTUserClaim struct {
    Login string
    jwt.StandardClaims
}

// Function to generating new token for login.
// Login must be autheticated or registered for generation
func GenerateJWTToken(claim *JWTUserClaim) (string, error) {
    // Get user session with personal user jwt key
    session, err := SessionStore.GetSession(claim.Login)

    if err != nil {
        return "", err
    }

    token := jwt.NewWithClaims(signMethod, claim)
    tokenString, _ := token.SignedString(session.JWTKey)

    return tokenString, nil
}

// Function to verify given token
// We must at first extract claim from unverified token to get token login
// Then, with extracted login, get appopriate jwt key and verity token
func VerifyToken(tokenString string) (*JWTUserClaim, bool) {
    claim := &JWTUserClaim{}
    _, parts, err := tokenParser.ParseUnverified(tokenString, claim)

    if err != nil {
        return claim, false
    }

    // Get session for parsed login
    session, err := SessionStore.GetSession(claim.Login)

    if err != nil {
        return claim, false
    }

    // Ignore alg header in token, use HMAC 256-bit always
    return claim, signMethod.Verify(strings.Join(parts[0:2], "."), parts[2], session.JWTKey) == nil
}
