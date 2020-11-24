package controllers

import (
    "net/http"
    "strings"
    "context"

    u "GoChessgameServer/util"
    "GoChessgameServer/store"

    jwt "github.com/dgrijalva/jwt-go"
)

// This function return a new handler,
// that will be check jwt tokens on valid
func TokenChecker(next http.Handler) http.Handler {

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        writeError := func(message string) {
            jsonResp := u.ErrorJson(message)
            w.WriteHeader(http.StatusForbidden)
            u.WriteResponse(w, jsonResp)
        }

        // Log requests
        contrLogger.Printf("TokenChecker: Request %s from %s\n", r.URL.Path, r.RemoteAddr)

        notAuthUris := []string {"/api/user/login", "/api/user/create"}
        urlPath := r.URL.Path

        // Check URIs, that is not needed in auth
        for _, uri := range notAuthUris {
            if uri == urlPath {
                next.ServeHTTP(w, r)
                return
            }
        }

        // Check Authorization header
        tokenHeader := r.Header.Get("Authorization")
        if tokenHeader == "" {
            writeError("Authorization header is missing")
            return
        }

        // Split then into Bearer XXXXXX (XXXXXX is jwt token)
        tokenSpl := strings.Split(tokenHeader, " ")
        if len(tokenSpl) != 2 && tokenSpl[0] != "Bearer" {
            writeError("Invalid/Mailformed token")
            return
        }

        // Get token and check signature
        token := tokenSpl[1]
        claims := &store.JWTClaims{}

        tokenRes, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
            return store.JWTKey, nil
        })

        if err != nil {
            writeError("Error when checking token signing")
            return
        }

        if !tokenRes.Valid {
            writeError("Token is not valid")
            return
        }

        contrLogger.Printf("TokenChecker: %s is making request", claims.Login)
        ctx := context.WithValue(r.Context(), "login", claims.Login)
        r = r.WithContext(ctx)
        next.ServeHTTP(w, r)
    })
}
