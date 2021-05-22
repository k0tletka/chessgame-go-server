package controllers

import (
    "net/http"
    "encoding/json"

    u "GoChessgameServer/util"
    "GoChessgameServer/database"

    "gorm.io/gorm"
)

// This function returns information about
// user account
func UserInfo(w http.ResponseWriter, r *http.Request) {

    writeError := u.WriteErrorCreator(w)
    login := r.Context().Value("login").(string)

    // Make database query
    var user database.User

    if err := database.DB.Find(&user, login).Error; err != nil {
        if err != gorm.ErrRecordNotFound {
            writeError("Connection error")
            w.WriteHeader(http.StatusInternalServerError)
            contrLogger.Printf("UserInfo: Error when executing query: %s\n", err.Error())
        } else {
            writeError("Oops, it seems that you account has been deleted. Please, restart you application")
        }

        return
    }

    // Send response
    resp := struct{
        Login string `json:"login"`
        Email string `json:"mail"`
        IsAdmin bool `json:"isadmin"`
    }{
        Login: user.Login,
        Email: user.Email,
        IsAdmin: user.IsAdmin,
    }

    if err := json.NewEncoder(w).Encode(resp); err != nil {
        writeError("Server error")
        w.WriteHeader(http.StatusInternalServerError)
        contrLogger.Printf("UserInfo: Error when sending response: %s\n", err.Error())
        return
    }
    w.Header().Add("Content-Type", "application/json")

    // Log new user
    contrLogger.Printf("UserInfo: User %s requested his info\n", user)
}
