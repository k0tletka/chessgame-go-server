package controllers

import (
    "net/http"
    "encoding/json"
    "reflect"

    u "GoChessgameServer/util"
    "GoChessgameServer/database"
    "GoChessgameServer/store"

    "golang.org/x/crypto/sha3"
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

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || reflect.DeepEqual(req, reqType{}) {
        writeError("Invalid request")
        contrLogger.Printf("LoginUsers: Error when parsing request from client: %s\n", err.Error())
        return
    }

    // Validate request
    if req.Login != "" && req.Password != "" {
        success := u.ValidateCredentials(req.Login, "", req.Password)
        if !success {
            writeError("Passed values is not matched with requirements to values")
            return
        }
    } else {
        writeError("One of the login or password is not specified, aborting")
        return
    }

    // Find user in the database
    results, err := database.QueryBlocking("SELECT TOP 1 * FROM dbo.Users WHERE Login = $1", req.Login)
    if err != nil {
        writeError("Connection error")
        w.WriteHeader(http.StatusInternalServerError)
        contrLogger.Printf("LoginUsers: Error when making query: %s\n", err.Error())
        return
    }
    if len(*results) == 0 {
        writeError("Login or password is not valid, please try again")
        return
    }

    // Calculating password hash and validating them
    user := (*results)[0]
    salt := user["PasswordHashSalt"].([]byte)
    realdigest := user["PasswordHash"].([]byte)

    digest1 := sha3.Sum256([]byte(req.Password))
    digest2 := sha3.Sum256(append(digest1[:], salt...))

    if string(digest2[:]) != string(realdigest) {
        writeError("Login or password is not valid, please try again")
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
