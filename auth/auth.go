package auth

import (
    "crypto/rand"
    "log"
    "os"

    "GoChessgameServer/logger"
    "GoChessgameServer/database"

    "golang.org/x/crypto/sha3"
    "gorm.io/gorm"
)

var (
    // Logger for module
    authLogger *log.Logger
)

func init() {
    // Create logger for auth module
    authLogger = logger.AddNewLogger("Auth", os.Stdout, log.LstdFlags | log.Lmsgprefix)
}

// This function performs login and password validating in the system
func AuthUser(login, password string) bool {
    // Check users in the database
    var user database.User

    if result := database.DB.Find(&user, login); result.Error != nil {
        if result.Error != gorm.ErrRecordNotFound {
            authLogger.Printf("Error when making query: %s\n", result.Error.Error())
        }

        return false
    }

    authed := checkPasswordValid(
        password,
        user.PasswordHash,
        user.PasswordHashSalt,
    )

    if !authed {
        authLogger.Printf("User %s has failed autheticating in system\n", login)
        return false
    }

    // Register session for logged user
    sinfo := &SessionInformation{
        IsAdmin: user.IsAdmin,
    }

    err := SessionStore.CreateNewSession(login, sinfo)
    return err == nil
}

// This function performs registering new users (session for new users also creates)
func RegisterUser(login, password, email string) bool {
    // Check if user axe exists already
    result := database.DB.Find(&database.User{}, login)

    if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
        authLogger.Printf("Error when making query: %s\n", result.Error.Error())
        return false
    }

    if result.Error == nil {
        return false
    }

    // Generate hash and salt for new account
    hash, salt := generateHashAndSalt(password)

    // Insert new account into db table
    result = database.DB.Create(&database.User{
        Login: login,
        Email: email,
        PasswordHash: hash,
        PasswordHashSalt: salt,
    })

    if result.Error != nil {
        authLogger.Printf("Error when executing query: %s\n", result.Error.Error())
        return false
    }

    // Register session for logged user
    sinfo := &SessionInformation{
        IsAdmin: false,
    }

    err := SessionStore.CreateNewSession(login, sinfo)
    return err == nil
}

// This function allows to change user password
func ChangeUserPassword(login, op, np string) bool {
    // Make request to get
    var user database.User

    if result := database.DB.Find(&user, login); result.Error != nil {
        if result.Error != gorm.ErrRecordNotFound {
            authLogger.Printf("Error when making query: %s\n", result.Error.Error())
        }

        return false
    }

    // Check old password valid
    if !checkPasswordValid(op, user.PasswordHash, user.PasswordHashSalt) {
        authLogger.Printf("User %s provided invalid password while trying to change it\n", login)
        return false
    }

    // Generate and update new password
    user.PasswordHash, user.PasswordHashSalt = generateHashAndSalt(np)

    if result := database.DB.Save(&user); result.Error != nil {
        authLogger.Printf("Error when executing query: %s\n", result.Error.Error())
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
