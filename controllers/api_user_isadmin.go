package controllers

import (
    "net/http"
    "encoding/json"

    u "GoChessgameServer/util"
    "GoChessgameServer/database"
)

// This function returns either true or false
// depending on if user is admin on this server
func IsAdmin(w http.ResponseWriter, r *http.Request) {

    writeError := func(message string) {
        jsonslice := u.ErrorJson(message)
        w.WriteHeader(http.StatusForbidden)
        u.WriteResponse(w, jsonslice)
    }

    // Get current user
    user := r.Context().Value("login").(string)

    // Make query to database
    results, err := database.QueryBlocking("SELECT TOP 1 IsAdmin FROM dbo.Users WHERE Login = $1", user)
    if err != nil {
        writeError("Connection error")
        w.WriteHeader(http.StatusInternalServerError)
        contrLogger.Printf("IsAdmin: Error when making request: %s\n", err.Error())
        return
    }
    if len(*results) == 0 {
        writeError("Oops, it seems that you account has been deleted. Please, restart you application")
        return
    }

    isAdmin := (*results)[0]["IsAdmin"].(bool)

    // Return isadmin to client
    resp := struct{
        IsAdmin bool `json:"isadmin"`
    }{
        IsAdmin: isAdmin,
    }

    if err = json.NewEncoder(w).Encode(resp); err != nil {
        writeError("Server error")
        w.WriteHeader(http.StatusInternalServerError)
        contrLogger.Printf("IsAdmin: Error when sending response: %s\n", err.Error())
        return
    }
    w.Header().Add("Content-Type", "application/json")

    // Log new user
    contrLogger.Printf("IsAdmin: User %s requested his admin status\n", user)
}
