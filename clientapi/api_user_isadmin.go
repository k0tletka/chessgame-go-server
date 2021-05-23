package clientapi

import (
    "net/http"
    "encoding/json"

    u "GoChessgameServer/util"
    "GoChessgameServer/auth"
)

// This function returns either true or false
// depending on if user is admin on this server
func IsAdmin(w http.ResponseWriter, r *http.Request) {

    writeError := u.WriteErrorCreator(w)

    // Get current user and his session
    user := r.Context().Value("login").(string)
    session, _ := auth.SessionStore.GetSession(user)

    // Return isadmin to client
    resp := struct{
        IsAdmin bool `json:"isadmin"`
    }{
        IsAdmin: session.IsAdmin,
    }

    if err := json.NewEncoder(w).Encode(resp); err != nil {
        writeError("Server error")
        w.WriteHeader(http.StatusInternalServerError)
        clientApiLogger.Printf("IsAdmin: Error when sending response: %s\n", err.Error())
        return
    }
    w.Header().Add("Content-Type", "application/json")

    // Log new user
    clientApiLogger.Printf("IsAdmin: User %s requested his admin status\n", user)
}
