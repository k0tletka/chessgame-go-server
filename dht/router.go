package dht

import (
    "sync"
    "net/http"
    "fmt"
    "time"

    c "GoChessgameServer/conf"
    u "GoChessgameServer/util"

    "github.com/gorilla/mux"
)

func InitializeDHTAPIServer(srvWaitor *sync.WaitGroup, srvResult chan<- *http.Server) {
    // Create new router for our api
    router := mux.NewRouter()

    // Handler websocket upgrade requests
    router.HandleFunc("/ws", DHTMgr.websocketHandler).Methods("GET")

    // Parse configuration values
    if !c.DecodeMetadata.IsDefined("dht_api", "use_tls") {
        dhtLogger.Fatalln("use_tls options is not defined, aborting to start DHT api server")
    }

    if !checkTLSMandatoryOptions() {
        dhtLogger.Fatalln("Needed config options for TLS is not defined, aborting to start DHT api server")
    }

    listenaddr, listenport := u.GetListenInformationServerAPI()

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
        dhtLogger.Fatalln(srv.ListenAndServeTLS(
            c.Conf.DHTApi.CertFile,
            c.Conf.DHTApi.KeyFile,
        ))
    } else {
        dhtLogger.Fatalln(srv.ListenAndServe())
    }
}

// Checks whenever cert and key file options
// defined in config file if TLS option enabled
func checkTLSMandatoryOptions() bool {
    if !c.Conf.DHTApi.UseTLS {
        return true
    }

    return c.DecodeMetadata.IsDefined("dht_api", "cert_file") &&
        c.DecodeMetadata.IsDefined("dht_api", "key_file")
}

// Shutdown HTTP server handler
func shutdownHTTPHandler() {
}
