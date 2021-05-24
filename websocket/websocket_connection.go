package websocket

import (
    "time"
    "sync"

    "github.com/gorilla/websocket"
)

var (
    // Ping sending duration for timer
    pingSendingDuration = 5 * time.Second
)

// This type represents websocket connection
type WebsocketConnection struct {
    conn *websocket.Conn

    // Pointer to store collection to have
    // ability to delete connection from list
    store *WebsocketStore

    // Pong channel
    pongChannel chan struct{}

    // Socket opened state and mutex for state
    openStateMutex *sync.Mutex
    openState bool

    // Handler function for read WS
    readHandlerFunction func(*WebsocketConnection, []byte)

    // List of handler on connection close
    connectionCloseHandlersMutex *sync.Mutex
    connectionCloseHandlers []func(*WebsocketConnection)

    // Channel for notifying that connection are closed
    ConnectionClosed chan bool
}

// Function to create new websocket connection
func NewWebsocketConnection(conn *websocket.Conn, readHandler func(*WebsocketConnection, []byte), store *WebsocketStore) *WebsocketConnection {
    pongChannel := make(chan struct{}, 1)

    result := &WebsocketConnection{
        conn: conn,
        store: store,
        pongChannel: pongChannel,
        openStateMutex: &sync.Mutex{},
        openState: true,
        readHandlerFunction: readHandler,
        connectionCloseHandlersMutex: &sync.Mutex{},
        connectionCloseHandlers: []func(*WebsocketConnection){},
        ConnectionClosed: make(chan bool, 1),
    }

    store.InsertConnection(result)

    // Set handler on pong messages
    conn.SetPongHandler(func (appData string) error {
        pongChannel <- struct{}{}
        return nil
    })

    go result.ReadGoroutine()
    go result.PingerGoroutine()
    return result
}

// Public methods
func (wc *WebsocketConnection) GetConnection() *websocket.Conn {
    return wc.conn
}

func (wc *WebsocketConnection) GetConnectionStore() *WebsocketStore {
    return wc.store
}

func (wc *WebsocketConnection) AddCloseHandler(handler func(*WebsocketConnection)) {
    wc.connectionCloseHandlersMutex.Lock()
    defer wc.connectionCloseHandlersMutex.Unlock()

    wc.connectionCloseHandlers = append(wc.connectionCloseHandlers, handler)
}

// Connection close method
func (wc *WebsocketConnection) CloseConnection(data string) {
    // Ignore errors
    wc.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, data))
    wc.closeConnectionForce()
}

func (wc *WebsocketConnection) closeConnectionForce() {
    wc.openStateMutex.Lock()
    defer wc.openStateMutex.Unlock()

    if wc.openState {
        // Close pong channel
        close(wc.pongChannel)

        // Close connection
        wc.conn.Close()
        wc.store.DeleteConnection(wc)

        for _, handler := range wc.connectionCloseHandlers {
            handler(wc)
        }

        wc.ConnectionClosed <- true

        wc.openState = false
    }
}

// Method for handling read messages from connection and calling readHandle function
func (wc *WebsocketConnection) ReadGoroutine() {
    for {
        mt, message, err := wc.conn.ReadMessage()

        if err != nil {
            wc.closeConnectionForce()
            return
        }

        if mt == websocket.TextMessage {
            // Call handler of text message
            go wc.readHandlerFunction(wc, message)
        }
    }
}

// Ping function that send ping request to peer
func (wc *WebsocketConnection) PingerGoroutine() {
    for {
        if err := wc.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(pingSendingDuration)); err != nil {
            wc.closeConnectionForce()
            return
        }

        timer := time.NewTimer(pingSendingDuration)

        select {
        case <-timer.C:
            websocketLogger.Printf("Websocket connection for %s expired, closing...\n", wc.conn.RemoteAddr())
            wc.closeConnectionForce()
            return
        case <-wc.pongChannel:
            time.Sleep(pingSendingDuration)
            continue
        }
    }
}
