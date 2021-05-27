package util

import (
    "encoding/json"
    "net/http"

    c "GoChessgameServer/conf"
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

// Function to get default listening port and addresses
func GetListenInformationClientAPI() (laddr string, lport uint16) {
    if !c.DecodeMetadata.IsDefined("client_api", "listenaddr") {
        laddr = "127.0.0.1"
    } else {
        laddr = c.Conf.CAPI.ListenAddr
    }

    if !c.DecodeMetadata.IsDefined("client_api", "listenport") {
        if c.Conf.CAPI.UseTLS {
            lport = 443
        } else {
            lport = 80
        }
    } else {
        lport = c.Conf.CAPI.ListenPort
    }

    return
}

func GetListenInformationGameAPI() (laddr string, lport uint16) {
    if !c.DecodeMetadata.IsDefined("game_api", "listenaddr") {
        laddr = "127.0.0.1"
    } else {
        laddr = c.Conf.GAPI.ListenAddr
    }

    if !c.DecodeMetadata.IsDefined("game_api", "listenport") {
        if c.Conf.GAPI.UseTLS {
            lport = 4443
        } else {
            lport = 800
        }
    } else {
        lport = c.Conf.GAPI.ListenPort
    }

    return
}

func GetListenInformationServerAPI() (laddr string, lport uint16) {
    if !c.DecodeMetadata.IsDefined("dht_api", "listenaddr") {
        laddr = "127.0.0.1"
    } else {
        laddr = c.Conf.DHTApi.ListenAddr
    }

    if !c.DecodeMetadata.IsDefined("dht_api", "listenport") {
        if c.Conf.DHTApi.UseTLS {
            lport = 4444
        } else {
            lport = 801
        }
    } else {
        lport = c.Conf.DHTApi.ListenPort
    }

    return
}

// Get public IP address
func GetPublicIPAddress() (string, error) {
    request, err := http.Get("http://ip-api.com/json")

    if err != nil {
        return "", err
    }

    response := struct{
        Query string
    }{}

    if err := json.NewDecoder(request.Body).Decode(&response); err != nil {
        return "", err
    }

    return response.Query, nil
}
