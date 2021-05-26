package clientapi

import (
    "net/http"
    "encoding/json"
    "fmt"

    u "GoChessgameServer/util"
    c "GoChessgameServer/conf"
    gameAPI "GoChessgameServer/gameapi"
)

func GetGameAPIUri(w http.ResponseWriter, r *http.Request) {

    writeError := u.WriteErrorCreator(w)
    listenaddr, listenport := gameAPI.GetListenInformation()

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
