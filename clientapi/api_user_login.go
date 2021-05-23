package clientapi

import (
    "net/http"
    "encoding/json"

    u "GoChessgameServer/util"
    "GoChessgameServer/auth"
)

// This controller perform user login
// into application
func LoginUsers(w http.ResponseWriter, r *http.Request) {

    writeError := u.WriteErrorCreator(w)

    type reqType struct {
        Login string `json:"login"`
        Password string `json:"pass"`
    }
    req := reqType{}

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError("Invalid request")
        clientApiLogger.Printf("LoginUsers: Error when parsing request from client: %s\n", err.Error())
        return
    }

    // Validate request
    success := u.ValidateValues(
        &u.VValue{Type: "Login", Value: req.Login},
        &u.VValue{Type: "Password", Value: req.Password},
    )

    if !success {
        writeError("Login or password doesn't satisfy value requirements")
        return
    }

    // Auth user
    if !auth.AuthUser(req.Login, req.Password) {
        writeError("Login or password incorrect, or user already logged in, please try again")
        return
    }

    // Password valid, generate jwt token
    claim := &auth.JWTUserClaim{Login: req.Login}
    token, err := auth.GenerateJWTToken(claim)

    if err != nil {
        writeError("Internal server error")
        return
    }
    // Write response
    resp := struct{
        Login string `json:"login"`
        JWTToken string `json:"token"`
    }{
        Login: req.Login,
        JWTToken: token,
    }

    if err = json.NewEncoder(w).Encode(resp); err != nil {
        writeError("Server error")
        w.WriteHeader(http.StatusInternalServerError)
        clientApiLogger.Printf("LoginUsers: Error when sending response: %s\n", err.Error())
        return
    }
    w.Header().Add("Content-Type", "application/json")

    // Log user logged in
    clientApiLogger.Printf("LoginUsers: User %s has been logged in\n", req.Login)
}
