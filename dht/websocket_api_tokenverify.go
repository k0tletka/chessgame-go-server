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
    var tokenByteData []byte
    var err error

    if tokenByteData, err = json.Marshal(&tokenData); err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Invalid request"))
        return
    }

    response := struct{
        Verified    bool            `json:"verified"`
        TokenData   json.RawMessage `json:"token_data"`
    }{
        Verified: verified,
        TokenData: tokenByteData,
    }

    var argsData []byte

    if argsData, err = json.Marshal(&response); err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Error when marshalling response: " + err.Error()))
        return
    }

    baseResponse := dhtAPIBaseRequest{
        MethodName: "tokenverify_response",
        Args: argsData,
    }

    var resData []byte

    if resData, err = json.Marshal(&baseResponse); err != nil {
        conn.WriteMessage(websocket.TextMessage, u.ErrorJson("Error when marshalling response: " + err.Error()))
        return
    }

    conn.WriteMessage(websocket.TextMessage, resData)
}
