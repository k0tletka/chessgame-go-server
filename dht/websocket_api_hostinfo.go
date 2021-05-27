package dht

import (
    "encoding/json"

    ws "GoChessgameServer/websocket"
    c "GoChessgameServer/util"
    u "GoChessgameServer/util"

    "github.com/gorilla/websocket"
)

func (m *DHTManager) hostinfoMethodHandler(wc *ws.WebsocketConnection, data *dhtAPIBaseRequest) {
    conn := wc.GetConnection()

    // Get listening info
    _, listenportCapi := c.GetListenInformationClientAPI()
    _, listenportGapi := c.GetListenInformationGameAPI()


    response := struct{
        ClientAPIPort   uint16  `json:"client_api_port"`
        GameAPIPort     uint16  `json:"game_api_port"`
    }{
        ClientAPIPort: listenportCapi,
        GameAPIPort: listenportGapi,
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
