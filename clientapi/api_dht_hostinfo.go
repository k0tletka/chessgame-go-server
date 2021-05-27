package clientapi

import (
    "net/http"
    "encoding/json"

    u "GoChessgameServer/util"
    "GoChessgameServer/dht"
)

func DhtHostInfo(w http.ResponseWriter, r *http.Request) {

    writeError := u.WriteErrorCreator(w)

    request := struct{
        ServerIdentifier    string  `json:"server_identifier"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        writeError("Invalid request")
        clientApiLogger.Printf("ChangePassword: Error when parsing request from client: %s\n", err.Error())
        return
    }

    hostInfo, err := dht.DHTMgr.GetHostInfoByServerIdentifier(request.ServerIdentifier)

    if err != nil {
        writeError("Can't get host info by server identifier: " + err.Error())
        return
    }

    // Encode result
    if err := json.NewEncoder(w).Encode(hostInfo); err != nil {
        writeError("Server error")
        w.WriteHeader(http.StatusInternalServerError)
        clientApiLogger.Printf("CreateLogin: Error when sending response: %s\n", err.Error())
        return
    }
    w.Header().Add("Content-Type", "application/json")
}
