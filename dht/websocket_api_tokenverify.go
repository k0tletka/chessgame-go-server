package dht

import (
    "encoding/json"

    ws "GoChessgameServer/websocket"
    u "GoChessgameServer/util"
    "GoChessgameServer/auth"

    "github.com/gorilla/websocket"
)

func (m *DHTManager) tokenverifyMethodHandler(wc *ws.WebsocketConnection, data *dhtAPIBaseRequest) {
    conn := wc.GetConnection()

    request := struct{
        TokenToVerify   string  `json:"verify"`
    }{}

    if err := json.Unmarshal(data.Args, &request); err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Invalid request"))
        return
    }

    // Verify given token
    tokenData, verified := auth.VerifyServerToken(request.TokenToVerify)

    // Write response
    var resData []byte
    var err error

    if resData, err = json.Marshal(&tokenData); err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Invalid request"))
        return
    }

    response := struct{
        Verified    bool            `json:"verified"`
        TokenData   json.RawMessage `json:"token_data"`
    }{
        Verified: verified,
        TokenData: resData,
    }

    if resData, err = json.Marshal(&response); err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Error when marshalling response: " + err.Error()))
        return
    }

    baseResponse := dhtAPIBaseRequest{
        MethodName: "tokenverify_response",
        Args: resData,
    }

    if resData, err = json.Marshal(&baseResponse); err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Error when marshalling response: " + err.Error()))
        return
    }

    conn.WriteMessage(websocket.TextMessage, resData)
}
