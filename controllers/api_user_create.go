package controllers

import (
    "net/http"
    "encoding/json"
    "regexp"
    "crypto/rand"

    u "GoChessgameServer/util"
    "GoChessgameServer/database"
    "GoChessgameServer/store"

    "golang.org/x/crypto/sha3"
    jwt "github.com/dgrijalva/jwt-go"
)

type newLogin struct {
    Login string `json:"login"`
    Email string `json:"mail"`
    Password string `json:"pass"`
}


// This controller create new users in the database
// and returns signed jwt token to client
func CreateLogin(w http.ResponseWriter, r *http.Request) {

    writeError := func(message string) {
        jsonslice := u.ErrorJson(message)
        w.WriteHeader(http.StatusForbidden)
        u.WriteResponse(w, jsonslice)
    }
    req := newLogin{}

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError("Invalid request")
        contrLogger.Printf("CreateLoggin: Error when parsing request from client: %s\n", err.Error())
        return
    }

    // Validate request
    success := validate(req)
    if !success {
        writeError("Passed values is not matched with requirements to values")
        return
    }

    // Check is logins and emails are existing
    results, err := database.QueryBlocking("SELECT Login, Email FROM dbo.Users WHERE Login = $1 OR Email = $2", req.Login, req.Email)
    if err != nil {
        writeError("Connection error")
        contrLogger.Printf("CreateLogin: Error when making query: %s\n", err.Error())
        return
    }
    if len(*results) > 0 {
        writeError("Account with specified login or email is existing already")
        return
    }

    // Hash password
    digest1 := sha3.Sum256([]byte(req.Password))

    // Generate salt, append and digest hash again
    salt := make([]byte, 256, 256)
    _, err = rand.Read(salt)
    if err != nil {
        writeError("Server error")
        w.WriteHeader(http.StatusInternalServerError)
        contrLogger.Printf("CreateLogin: Error when generating salt: %s\n", err.Error())
        return
    }
    digest2 := sha3.Sum256(append(digest1[:], salt...))

    // Insert new account into db table
    _, err = database.QueryExecBlocking(`
    INSERT INTO dbo.Users(Login, Email, PasswordHash, PasswordHashSalt)
    VALUES ($1, $2, $3, $4)`, req.Login, req.Email, string(digest2[:]), string(salt))
    if err != nil {
        writeError("Server error")
        w.WriteHeader(http.StatusInternalServerError)
        contrLogger.Printf("CreateLogin: Error when executing query: %s\n", err.Error())
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

func validate(login newLogin) bool {

    // Validate login (must be with lower symbols and numbers, and from 6 to 100 symbols)
    if matched, err := regexp.MatchString(`^[a-z0-9]{6,100}$`, login.Login); err != nil || !matched {
        return false
    }

    // Validate email
    if matched, err := regexp.MatchString(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)+$`, login.Email); err != nil || !matched {
        return false
    }

    // Validate password
    if matched, err := regexp.MatchString(`[0-9]`, login.Password); err != nil || !matched {
        return false
    }
    if matched, err := regexp.MatchString(`[A-Z]`, login.Password); err != nil || !matched {
        return false
    }
    if matched, err := regexp.MatchString(`^\S{8,}$`, login.Password); err != nil || !matched {
        return false
    }

    return true
}
