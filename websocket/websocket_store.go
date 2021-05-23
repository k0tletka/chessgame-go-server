package websocket

import (
    "sync"
)

// This type represents store object for websocket connection list
type WebsocketStore struct {
    sync.Mutex
    wsConns map[*WebsocketConnection]struct{}
}

func NewWebsocketStore() *WebsocketStore {
    return &WebsocketStore{
        wsConns: make(map[*WebsocketConnection]struct{}),
    }
}

func (w *WebsocketStore) InsertConnection(conn *WebsocketConnection) {
    w.Lock()
    w.wsConns[conn] = struct{}{}
    w.Unlock()
}

func (w *WebsocketStore) DeleteConnection(conn *WebsocketConnection) {
    w.Lock()
    delete(w.wsConns, conn)
    w.Unlock()
}

func (w *WebsocketStore) GetConnections() []*WebsocketConnection {
    result := []*WebsocketConnection{}

    for k, _ := range w.wsConns {
        result = append(result, k)
    }

    return result
}
