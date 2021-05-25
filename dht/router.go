package dht

import (
    "sync"
    "net/http"
    "fmt"
    "time"

    c "GoChessgameServer/conf"

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

    listenaddr, listenport := getListenInformation()

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

// Gets address and port for server API listening
func getListenInformation() (laddr string, lport uint16) {
    if !c.DecodeMetadata.IsDefined("dht_api", "listenaddr") {
        laddr = "127.0.0.1"
    } else {
        laddr = c.Conf.DHTApi.ListenAddr
    }

    if !c.DecodeMetadata.IsDefined("dht_api", "listenport") {
        if c.Conf.DHTApi.UseTLS {
            lport = 4444
        } else {
            lport = 801
        }
    } else {
        lport = c.Conf.DHTApi.ListenPort
    }

    return
}

// Shutdown HTTP server handler
func shutdownHTTPHandler() {
}
