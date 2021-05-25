package clientapi

import (
    "fmt"
    "net/http"
    "time"
    "sync"

    c "GoChessgameServer/conf"

    "github.com/gorilla/mux"
)

// Function to initialize router for client API
func InitializeClientAPIServer(srvWaitor *sync.WaitGroup, srvResult chan<- *http.Server) {
    // Create our API router
    router := mux.NewRouter()

    // Handlers on a variety of controllers
    router.HandleFunc("/api/user/create", CreateLogin).Methods("POST")
    router.HandleFunc("/api/user/login", LoginUsers).Methods("POST")
    router.HandleFunc("/api/user/disconnect", DisconnectUser).Methods("GET")
    router.HandleFunc("/api/user/isadmin", IsAdmin).Methods("GET")
    router.HandleFunc("/api/user/changepass", ChangePassword).Methods("POST")
    router.HandleFunc("/api/user/info", UserInfo).Methods("GET")
    router.HandleFunc("/api/playerstat", UserStatistic).Methods("GET", "POST")
    router.HandleFunc("/api/motd", GetMotd).Methods("GET")

    // Middleware for token checking
    router.Use(TokenChecker)

    // Setup http handler options and start client API server
    if !c.DecodeMetadata.IsDefined("client_api", "use_tls") {
        clientApiLogger.Fatalln("use_tls options is not defined, aborting to start client api server")
    }

    if !checkTLSMandatoryOptions() {
        clientApiLogger.Fatalln("Needed config options for TLS not defined, aborting to start client api server")
    }

    listenaddr, listenport := getListenInformation()

    srv := &http.Server{
        Handler: router,
        Addr: fmt.Sprintf("%s:%d", listenaddr, listenport),
        WriteTimeout: 15 * time.Second,
        ReadTimeout: 15 * time.Second,
    }

    // Norify main thread and return srv
    srvResult <- srv
    srvWaitor.Done()

    if c.Conf.CAPI.UseTLS {
        clientApiLogger.Fatalln(srv.ListenAndServeTLS(
            c.Conf.CAPI.CertFile,
            c.Conf.CAPI.KeyFile,
        ))
    } else {
        clientApiLogger.Fatalln(srv.ListenAndServe())
    }
}

// Checks whenever cert and key file options
// defined in config file if TLS option enabled
func checkTLSMandatoryOptions() bool {
    if !c.Conf.CAPI.UseTLS {
        return true
    }

    return c.DecodeMetadata.IsDefined("client_api", "cert_file") &&
        c.DecodeMetadata.IsDefined("client_api", "key_file")
}


// Gets address and port for server API listening
func getListenInformation() (laddr string, lport uint16) {
    if !c.DecodeMetadata.IsDefined("client_api", "listenaddr") {
        laddr = "127.0.0.1"
    } else {
        laddr = c.Conf.CAPI.ListenAddr
    }

    if !c.DecodeMetadata.IsDefined("client_api", "listenport") {
        if c.Conf.CAPI.UseTLS {
            lport = 443
        } else {
            lport = 80
        }
    } else {
        lport = c.Conf.CAPI.ListenPort
    }

    return
}
