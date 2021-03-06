package dht

import (
    ws "GoChessgameServer/websocket"
)

func (m *DHTManager) hostinfoResponseMethodHandler(wc *ws.WebsocketConnection, data *dhtAPIBaseRequest) {
    // Return request to caller code
    if ch, ok := m.readSynchronizationChannels[wc]; ok {
        ch <- data
    }
}
