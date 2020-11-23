package util

import (
    "encoding/json"
    "net/http"
)

func ErrorJson(msgid int, message string) ([]byte, error) {
    return json.Marshal(map[string]interface{} { "msgid": msgid, "error": message })
}

func MessageJson(msgid int, message string) ([]byte, error) {
    return json.Marshal(map[string]interface{} { "msgid": msgid, "message": message })
}

func WriteResponse(w http.ResponseWriter, jsonSlice []byte) {
    w.Header().Add("Content-Type", "application/json")
    _, _ = w.Write(jsonSlice)
}
