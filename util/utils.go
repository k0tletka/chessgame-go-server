package util

import (
    "encoding/json"
    "net/http"
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
