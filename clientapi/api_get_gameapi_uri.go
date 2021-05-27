package clientapi

import (
    "net/http"
    "encoding/json"
    "fmt"

    u "GoChessgameServer/util"
    c "GoChessgameServer/conf"
)

func GetGameAPIUri(w http.ResponseWriter, r *http.Request) {

    writeError := u.WriteErrorCreator(w)
    listenaddr, listenport := u.GetListenInformationGameAPI()

    if listenaddr == "0.0.0.0" {
        var err error
        listenaddr, err = u.GetPublicIPAddress()

        if err != nil {
            writeError("Error when getting public ip address")
            return
        }
    }

    response := struct{
        GameAPIEndpoint string  `json:"gameapi_endpoint"`
    }{
        GameAPIEndpoint: fmt.Sprintf(
            "ws%s://%s:%d/ws",
            map[bool]string{true: "s", false: ""}[c.Conf.GAPI.UseTLS],
            listenaddr,
            listenport,
        ),
    }

    if err := json.NewEncoder(w).Encode(&response); err != nil {
        writeError("Server error")
        w.WriteHeader(http.StatusInternalServerError)
        clientApiLogger.Printf("GameAPIEndpoint: Error when sending response: %s\n", err.Error())
        return
    }
    w.Header().Add("Content-Type", "application/json")
}
