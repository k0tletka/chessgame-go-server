package clientapi

import (
    "net/http"
    "io"

    "GoChessgameServer/store"
)

// This function returns stored motd
// to users
func GetMotd(w http.ResponseWriter, r *http.Request) {

    contextUser := r.Context().Value("login").(string)

    // Set content type to text/markdown and return this string
    w.Header().Set("Content-Type", "text/markdown")
    w.WriteHeader(http.StatusOK)

    io.WriteString(w, store.MotdString)

    // Log
    clientApiLogger.Printf("GetMotd: User %s requested server motd\n", contextUser)
}
