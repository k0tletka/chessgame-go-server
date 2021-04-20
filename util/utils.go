package util

import (
    "encoding/json"
    "net/http"
)

// Utility functions
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

// This function returns writeError function with specified ResponseWriter variable
func WriteErrorCreator(w http.ResponseWriter) func(string) {

    return func(message string) {
        jsonslice := ErrorJson(message)
        w.WriteHeader(http.StatusForbidden)
        WriteResponse(w, jsonslice)
    }
}

// Function that performs abs with int valies
func Abs(n int) int {
    y := n >> 31
    return (y ^ n) - y
}
