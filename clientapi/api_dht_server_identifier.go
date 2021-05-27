package clientapi

import (
    "net/http"
    "encoding/json"

    "GoChessgameServer/dht"
    u "GoChessgameServer/util"
)

// This handler returns server identifier of this instance to client
func DhtServerIdentifier(w http.ResponseWriter, r *http.Request) {

    writeError := u.WriteErrorCreator(w)

    response := struct{
        ServerIdentifier    string  `json:"server_identifier"`
    }{
        ServerIdentifier: dht.DHTMgr.GetServerIdentifier(),
    }

    if err := json.NewEncoder(w).Encode(&response); err != nil {
        writeError("Server error")
        w.WriteHeader(http.StatusInternalServerError)
        clientApiLogger.Printf("ServerIdentifier: Error when sending response: %s\n", err.Error())
        return
    }
    w.Header().Add("Content-Type", "application/json")
}
