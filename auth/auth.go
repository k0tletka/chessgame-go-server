package auth

import (
    "crypto/rand"
    "log"
    "os"

    "GoChessgameServer/logger"
    "GoChessgameServer/database"

    "golang.org/x/crypto/sha3"
)

// Logger for module
var authLogger *log.Logger

func init() {
    // Create logger for auth module
    authLogger = logger.AddNewLogger("Auth", os.Stdout, log.LstdFlags | log.Lmsgprefix)
}

// This function performs login and password validating in the system
func AuthUser(login, password string) bool {
    // Check users in the database
    results, err := database.QueryBlocking("SELECT TOP 1 * FROM dbo.Users WHERE Login = $1", login)

    if err != nil {
        authLogger.Printf("Error when making query: %s\n", err.Error())
        return false
    }

    if len(*results) == 0 {
        return false
    }

    // User info
    user := (*results)[0]
    salt := user["PasswordHashSalt"].([]byte)
    hash := user["PasswordHash"].([]byte)
    isAdmin := user["IsAdmin"].(bool)

    authed := checkPasswordValid(password, hash, salt)

    if !authed {
        authLogger.Printf("User %s has failed autheticating in system\n", login)
        return false
    }

    // Register session for logged user
    sinfo := &SessionInformation{
        IsAdmin: isAdmin,
    }

    err = SessionStore.CreateNewSession(login, sinfo)
    return err == nil
}

// This function performs registering new users (session for new users also creates)
func RegisterUser(login, password, email string) bool {
    // Check if user axe exists already
    results, err := database.QueryBlocking("SELECT Login, Email FROM dbo.Users WHERE Login = $1", login)

    if err != nil {
        authLogger.Printf("Error when making query: %s\n", err.Error())
        return false
    }

    if len(*results) > 0 {
        return false
    }

    // Generate hash and salt for new account
    hash, salt := generateHashAndSalt(password)

    // Insert new account into db table
    _, err = database.QueryExecBlocking(`
    INSERT INTO dbo.Users(Login, Email, PasswordHash, PasswordHashSalt)
    VALUES ($1, $2, $3, $4)`, login, email, hash, salt)

    if err != nil {
        authLogger.Printf("Error when executing query: %s\n", err.Error())
        return false
    }

    // Register session for logged user
    sinfo := &SessionInformation{
        IsAdmin: false,
    }

    err = SessionStore.CreateNewSession(login, sinfo)
    return err == nil
}

// This function allows to change user password
func ChangeUserPassword(login, op, np string) bool {
    // Make request to get
    results, err := database.QueryBlocking(`SELECT PasswordHash, PasswordHashSalt FROM dbo.Users WHERE Login = $1`, login)

    if err != nil {
        authLogger.Printf("Error when executing query: %s\n", err.Error())
        return false
    }

    if len(*results) == 0 {
        return false
    }

    // Check old password valid
    opHash := (*results)[0]["PasswordHash"].([]byte)
    opHashSalt := (*results)[0]["PasswordHashSalt"].([]byte)

    if !checkPasswordValid(op, opHash, opHashSalt) {
        authLogger.Printf("User %s provided invalid password while trying to change it\n", login)
        return false
    }

    // Generate and update new password
    hash, salt := generateHashAndSalt(np)

    _, err = database.QueryExecBlocking(`
    UPDATE dbo.Users SET PasswordHash = $1, PasswordHashSalt = $2
    FROM dbo.Users
    WHERE Login = $3`, hash, salt, login)

    if err != nil {
        authLogger.Printf("Error when executing query: %s\n", err.Error())
        return false
    }

    return true
}

// Utility functions

// This function check user password with given hash and salt.
// Return true if password valid for given hash, otherwise false
func checkPasswordValid(password string, hash, salt []byte) bool {
    digest1 := sha3.Sum256([]byte(password))
    digest2 := sha3.Sum256(append(digest1[:], salt...))

    return string(digest2[:]) == string(hash)
}

// This function generates new hash and salt for given password
func generateHashAndSalt(password string) (hash, salt []byte) {
    digest1 := sha3.Sum256([]byte(password))

    // Generate salt, append and digest hash again
    salt = make([]byte, 256, 256)
    n, err := rand.Read(salt)

    if err != nil {
        // Fill missing bytes with zeros
        for i := n; i < n; i++ {
            salt[i] = byte(0)
        }
    }

    digest2 := sha3.Sum256(append(digest1[:], salt...))

    return digest2[:], salt
}
