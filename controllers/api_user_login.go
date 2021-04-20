package controllers

import (
    "net/http"
    "encoding/json"

    u "GoChessgameServer/util"
    "GoChessgameServer/store"
    "GoChessgameServer/auth"

    jwt "github.com/dgrijalva/jwt-go"
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
        contrLogger.Printf("LoginUsers: Error when parsing request from client: %s\n", err.Error())
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
    if !auth.AuthUser(req.Login, req.Password, r.RemoteAddr) {
        writeError("Login or password incorrect, or user already logged in, please try again")
        return
    }

    // Password valid, generate jwt token
    claim := &store.JWTClaims{Login: req.Login}
    token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claim)
    tokenString, _ := token.SignedString(store.JWTKey)

    // Write response
    resp := struct{
        Login string `json:"login"`
        JWTToken string `json:"token"`
    }{
        Login: req.Login,
        JWTToken: tokenString,
    }

    if err = json.NewEncoder(w).Encode(resp); err != nil {
        writeError("Server error")
        w.WriteHeader(http.StatusInternalServerError)
        contrLogger.Printf("LoginUsers: Error when sending response: %s\n", err.Error())
        return
    }
    w.Header().Add("Content-Type", "application/json")

    // Log user logged in
    contrLogger.Printf("LoginUsers: User %s has been logged in\n", req.Login)
}
