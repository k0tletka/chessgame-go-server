package controllers

import (
    "net/http"
    "encoding/json"

    u "GoChessgameServer/util"
    "GoChessgameServer/store"

    jwt "github.com/dgrijalva/jwt-go"
)

// This controller create new users in the database
// and returns signed jwt token to client
func CreateLogin(w http.ResponseWriter, r *http.Request) {

    writeError := u.WriteErrorCreator(w)

    type reqType struct {
        Login string `json:"login"`
        Email string `json:"mail"`
        Password string `json:"pass"`
    }
    req := reqType{}

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError("Invalid request")
        contrLogger.Printf("CreateLoggin: Error when parsing request from client: %s\n", err.Error())
        return
    }

    // Validate request
    success := u.ValidateValues(
        &u.VValue{Type: "Login", Value: req.Login},
        &u.VValue{Type: "Password", Value: req.Password},
        &u.VValue{Type: "Email", Value: req.Email},
    )

    if !success {
        writeError("Login, password or email doesn't satisfy value requirements")
        return
    }

    if !auth.RegisterUser(req.Login, req.Password, req.Email, r.RemoteAddr) {
        writeError("Register failed: maybe, login is occupied by another account or internal error occured")
        return
    }

    // Create new token
    claim := &store.JWTClaims{Login: req.Login}
    token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claim)
    tokenString, _ := token.SignedString(store.JWTKey)

    // Return to client login and his jwt token
    resp := struct {
        Login string `json:"login"`
        JWTToken string `json:"token"`
    }{
        Login: req.Login,
        JWTToken: tokenString,
    }

    if err = json.NewEncoder(w).Encode(resp); err != nil {
        writeError("Server error")
        w.WriteHeader(http.StatusInternalServerError)
        contrLogger.Printf("CreateLogin: Error when sending response: %s\n", err.Error())
        return
    }
    w.Header().Add("Content-Type", "application/json")

    // Log new user
    contrLogger.Printf("CreateLogin: Created new user %s\n", req.Login)
}
