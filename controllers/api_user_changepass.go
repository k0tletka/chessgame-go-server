package controllers

import (
    "net/http"
    "encoding/json"

    u "GoChessgameServer/util"
    "GoChessgameServer/auth"
)

// This function implements password changing
func ChangePassword(w http.ResponseWriter, r *http.Request) {

    writeError := u.WriteErrorCreator(w)
    user := r.Context().Value("login").(string)

    type reqType struct {
        OldPassword string `json:"oldpass"`
        NewPassword string `json:"newpass"`
    }
    req := reqType{}

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError("Invalid request")
        contrLogger.Printf("ChangePassword: Error when parsing request from client: %s\n", err.Error())
        return
    }

    // Validate passed password
    success := u.ValidateValues(&u.VValue{Type: "Password", Value: req.NewPassword})

    if !success {
        writeError("Password format doesn't satisfy security requirements")
        return
    }

    if !auth.ChangeUserPassword(user, req.OldPassword, req.NewPassword) {
        writeError("Old password is invalid or internal error occured")
        return
    }

    // Log password changed
    contrLogger.Printf("ChangePassword: Changed password for user %s\n", user)
}
