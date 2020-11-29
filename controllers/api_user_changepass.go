package controllers

import (
    "net/http"
    "encoding/json"
    "crypto/rand"
    "reflect"

    u "GoChessgameServer/util"
    "GoChessgameServer/database"

    "golang.org/x/crypto/sha3"
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

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || reflect.DeepEqual(req, reqType{}) {
        writeError("Invalid request")
        contrLogger.Printf("ChangePassword: Error when parsing request from client: %s\n", err.Error())
        return
    }

    // Validate passed password
    if req.OldPassword == "" || req.NewPassword == "" || !u.ValidateCredentials("", "", req.OldPassword) || !u.ValidateCredentials("", "", req.NewPassword) {
        writeError("New password is not corresponds to password requirements")
        return
    }

    if req.OldPassword == req.NewPassword {
        writeError("New password is the same as old password, aborting")
        return
    }

    // Check password if its not the same password as old one
    results, err := database.QueryBlocking(`SELECT PasswordHash, PasswordHashSalt FROM dbo.Users WHERE Login = %1`, user)
    if err != nil {
        writeError("Connection error")
        w.WriteHeader(http.StatusInternalServerError)
        contrLogger.Printf("ChangePassword: Error when executing query: %s\n", err.Error())
        return
    }
    if len(*results) == 0 {
        writeError("Oops, it seems that you account has been deleted. Please, restart you application")
        return
    }

    oldPasswordHash := (*results)[0]["PasswordHash"].([]byte)
    oldPasswordHashSalt := (*results)[0]["PasswordHashSalt"].([]byte)

    // Check old password
    oldPasswordDigest1 := sha3.Sum256([]byte(req.OldPassword))
    oldPasswordDigest2 := sha3.Sum256(append(oldPasswordDigest1[:], oldPasswordHashSalt...))

    if string(oldPasswordDigest2[:]) != string(oldPasswordHash) {
        writeError("Old password is not valid")
        return
    }

    // Generate new password
    newPasswordDigest1 := sha3.Sum256([]byte(req.NewPassword))
    newPasswordSalt := make([]byte, 256, 256)
    _, err = rand.Read(newPasswordSalt)
    if err != nil {
        writeError("Server error")
        w.WriteHeader(http.StatusInternalServerError)
        contrLogger.Printf("ChangePassword: Error when generating salt: %s\n", err.Error())
        return
    }
    newPasswordDigest2 := sha3.Sum256(append(newPasswordDigest1[:], newPasswordSalt...))

    // Update password
    result, err := database.QueryExecBlocking(`
    UPDATE dbo.Users SET PasswordHash = %1, PasswordHashSalt = %2
    FROM dbo.Users
    WHERE Login = %1`, newPasswordDigest2[:], newPasswordSalt)
    if err != nil {
        writeError("Connection error")
        w.WriteHeader(http.StatusInternalServerError)
        contrLogger.Printf("ChangePassword: Error when executing query: %s\n", err.Error())
        return
    }

    rowsAffected, _ := (*result).RowsAffected()
    if rowsAffected < int64(1) {
        writeError("Oops, it seems that you account has been deleted. Please, restart you application")
        return
    }

    // Log password changed
    contrLogger.Printf("ChangePassword: Changed password for user %s\n", user)
}
