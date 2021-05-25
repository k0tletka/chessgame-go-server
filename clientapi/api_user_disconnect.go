package clientapi

import (
    "net/http"

    u "GoChessgameServer/util"
    "GoChessgameServer/auth"
)

// Function for deleting user from session list
func DisconnectUser(w http.ResponseWriter, r *http.Request) {

    writeError := u.WriteErrorCreator(w)
    login := r.Context().Value("login").(string)

    // Delete session from list
    err := auth.SessionStore.TerminateSession(login)
    if err != nil {
        writeError("Session with that login not found")
        return
    }
}
