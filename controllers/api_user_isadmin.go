package controllers

import (
    "net/http"
    "encoding/json"

    u "GoChessgameServer/util"
)

// This function returns either true or false
// depending on if user is admin on this server
func IsAdmin(w http.ResponseWriter, r *http.Request) {

    writeError := u.WriteErrorCreator(w)

    // Get current user
    user := r.Context().Value("login").(string)
    isAdmin := r.Context().Value("isadmin").(bool)

    // Return isadmin to client
    resp := struct{
        IsAdmin bool `json:"isadmin"`
    }{
        IsAdmin: isAdmin,
    }

    if err := json.NewEncoder(w).Encode(resp); err != nil {
        writeError("Server error")
        w.WriteHeader(http.StatusInternalServerError)
        contrLogger.Printf("IsAdmin: Error when sending response: %s\n", err.Error())
        return
    }
    w.Header().Add("Content-Type", "application/json")

    // Log new user
    contrLogger.Printf("IsAdmin: User %s requested his admin status\n", user)
}
