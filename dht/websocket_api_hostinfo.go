package dht

import (
    "encoding/json"

    ws "GoChessgameServer/websocket"
    c "GoChessgameServer/conf"
    u "GoChessgameServer/util"

    "github.com/gorilla/websocket"
)

func (m *DHTManager) hostinfoMethodHandler(wc *ws.WebsocketConnection, data *dhtAPIBaseRequest) {
    conn := wc.GetConnection()

    // Get listening info
    _, listenportCapi := u.GetListenInformationClientAPI()
    _, listenportGapi := u.GetListenInformationGameAPI()


    response := struct{
        ClientAPIPort   uint16  `json:"client_api_port"`
        ClientAPITLS    bool    `json:"client_api_tls"`
        GameAPIPort     uint16  `json:"game_api_port"`
        GameAPITLS      bool    `json:"game_api_tls"`
    }{
        ClientAPIPort: listenportCapi,
        ClientAPITLS: c.Conf.CAPI.UseTLS,
        GameAPIPort: listenportGapi,
        GameAPITLS: c.Conf.GAPI.UseTLS,
    }

    var resData []byte
    var err error

    if resData, err = json.Marshal(&response); err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Error occured when marshalling request: " + err.Error()))
        return
    }

    baseResponse := dhtAPIBaseRequest{
        MethodName: "hostinfo_response",
        Args: resData,
    }

    if resData, err = json.Marshal(&baseResponse); err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Error occured when marshalling request: " + err.Error()))
        return
    } else {
        conn.WriteMessage(websocket.TextMessage, resData)
    }
}
