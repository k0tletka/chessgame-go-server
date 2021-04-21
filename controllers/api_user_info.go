package controllers

import (
    "net/http"
    "encoding/json"

    u "GoChessgameServer/util"
    "GoChessgameServer/database"
)

// This function returns information about
// user account
func UserInfo(w http.ResponseWriter, r *http.Request) {

    writeError := u.WriteErrorCreator(w)
    user := r.Context().Value("login").(string)

    // Make database query
    results, err := database.QueryBlocking(`
    SELECT Email, IsAdmin
    FROM dbo.Users
    WHERE Login = $1`, user)

    if err != nil {
        writeError("Connection error")
        w.WriteHeader(http.StatusInternalServerError)
        contrLogger.Printf("UserInfo: Error when executing query: %s\n", err.Error())
        return
    }

    if len(*results) == 0 {
        writeError("Oops, it seems that you account has been deleted. Please, restart you application")
        return
    }

    email := (*results)[0]["Email"].(string)
    isAdmin := (*results)[0]["IsAdmin"].(bool)

    // Send response
    resp := struct{
        Login string `json:"login"`
        Email string `json:"mail"`
        IsAdmin bool `json:"isadmin"`
    }{
        Login: user,
        Email: email,
        IsAdmin: isAdmin,
    }

    if err = json.NewEncoder(w).Encode(resp); err != nil {
        writeError("Server error")
        w.WriteHeader(http.StatusInternalServerError)
        contrLogger.Printf("UserInfo: Error when sending response: %s\n", err.Error())
        return
    }
    w.Header().Add("Content-Type", "application/json")

    // Log new user
    contrLogger.Printf("UserInfo: User %s requested his info\n", user)
}
