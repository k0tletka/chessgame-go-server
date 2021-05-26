package clientapi

import (
    "net/http"
    "strings"
    "context"

    u "GoChessgameServer/util"
    "GoChessgameServer/auth"
)

// This function return a new handler,
// that will be check jwt tokens on valid
func TokenChecker(next http.Handler) http.Handler {

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        writeError := u.WriteErrorCreator(w)
        notAuthUris := []string {"/api/user/login", "/api/user/create", "/api/gameapi_endpoint"}
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
        claim, verified := auth.VerifyToken(token)

        if !verified {
            writeError("Unautheticated request or maybe session has been terminated")
            return
        }

        ctx := context.WithValue(r.Context(), "login", claim.Login)
        r = r.WithContext(ctx)
        next.ServeHTTP(w, r)
    })
}
