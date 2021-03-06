package gameapi

import (
    "net/http"
    "sync"
    "time"
    "fmt"

    c "GoChessgameServer/conf"
    u "GoChessgameServer/util"
    ws "GoChessgameServer/websocket"

    "github.com/gorilla/mux"
)

var (
    // Websocket store object
    wsStore = ws.NewWebsocketStore()
)

// Function to initialize router for game API
func InitializeGameAPIServer(srvWaitor *sync.WaitGroup, srvResult chan<- *http.Server) {
    // Create our API router
    router := mux.NewRouter()

    // Handle requests to websocket init function
    router.HandleFunc("/ws", WebsocketHandler).Methods("GET")

    // Parse configuration values
    if !c.DecodeMetadata.IsDefined("game_api", "use_tls") {
        gameApiLogger.Fatalln("use_tls options is not defined, aborting to start game api server")
    }

    if !checkTLSMandatoryOptions() {
        gameApiLogger.Fatalln("Needed config options for TLS is not defined, aborting to start game api server")
    }

    listenaddr, listenport := u.GetListenInformationGameAPI()

    srv := &http.Server{
        Handler: router,
        Addr: fmt.Sprintf("%s:%d", listenaddr, listenport),
        WriteTimeout: 15 * time.Second,
        ReadTimeout: 15 * time.Second,
    }

    // Handler function on shutdown to close all WS connections
    srv.RegisterOnShutdown(shutdownHTTPHandler)

    // Notify main thread and return srv
    srvResult <- srv
    srvWaitor.Done()

    if c.Conf.GAPI.UseTLS {
        gameApiLogger.Fatalln(srv.ListenAndServeTLS(
            c.Conf.GAPI.CertFile,
            c.Conf.GAPI.KeyFile,
        ))
    } else {
        gameApiLogger.Fatalln(srv.ListenAndServe())
    }
}

// Checks whenever cert and key file options
// defined in config file if TLS option enabled
func checkTLSMandatoryOptions() bool {
    if !c.Conf.GAPI.UseTLS {
        return true
    }

    return c.DecodeMetadata.IsDefined("game_api", "cert_file") &&
        c.DecodeMetadata.IsDefined("game_api", "key_file")
}

// Shutdown HTTP server handler
func shutdownHTTPHandler() {
    // Iterate over WS connections and close all
    for _, conn := range wsStore.GetConnections() {
        conn.CloseConnection("Server is shutdowning.")
    }
}
