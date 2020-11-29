package util

import (
    "encoding/json"
    "net/http"
    "regexp"
)

func ErrorJson(message string) []byte {
    res, _ := json.Marshal(map[string]interface{} { "error": message })
    return res
}

func MessageJson(message string) []byte {
    res, _ := json.Marshal(map[string]interface{} { "message": message })
    return res
}

func WriteResponse(w http.ResponseWriter, jsonSlice []byte) {
    w.Header().Add("Content-Type", "application/json")
    _, _ = w.Write(jsonSlice)
}

// This function validates can validate login, email and password
func ValidateCredentials(login string, email string, password string) bool {

    // Validate login (must be with lower symbols and numbers, and from 6 to 100 symbols)
    if login != "" {
        if matched, err := regexp.MatchString(`^[a-z0-9]{6,100}$`, login); err != nil || !matched {
            return false
        }
    }

    // Validate email
    if email != "" {
        if matched, err := regexp.MatchString(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)+$`, email); err != nil || !matched {
            return false
        }
    }

    // Validate password
    if password != "" {
        if matched, err := regexp.MatchString(`[0-9]`, password); err != nil || !matched {
            return false
        }
        if matched, err := regexp.MatchString(`[A-Z]`, password); err != nil || !matched {
            return false
        }
        if matched, err := regexp.MatchString(`^\S{8,}$`, password); err != nil || !matched {
            return false
        }
    }

    return true
}

// This function returns writeError function with specified ResponseWriter variable
func WriteErrorCreator(w http.ResponseWriter) func(string) {

    return func(message string) {
        jsonslice := ErrorJson(message)
        w.WriteHeader(http.StatusForbidden)
        WriteResponse(w, jsonslice)
    }
}
