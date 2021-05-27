package clientapi

import (
    "fmt"
    "net/http"
    "time"
    "sync"

    c "GoChessgameServer/conf"
    u "GoChessgameServer/util"

    "github.com/gorilla/mux"
)

// Function to initialize router for client API
func InitializeClientAPIServer(srvWaitor *sync.WaitGroup, srvResult chan<- *http.Server) {
    // Create our API router
    router := mux.NewRouter()

    // Handlers on a variety of controllers
    router.HandleFunc("/api/lobby/list", LobbyList).Methods("GET")
    router.HandleFunc("/api/user/create", CreateLogin).Methods("POST")
    router.HandleFunc("/api/user/login", LoginUsers).Methods("POST")
    router.HandleFunc("/api/user/disconnect", DisconnectUser).Methods("GET")
    router.HandleFunc("/api/user/isadmin", IsAdmin).Methods("GET")
    router.HandleFunc("/api/user/changepass", ChangePassword).Methods("POST")
    router.HandleFunc("/api/user/info", UserInfo).Methods("GET")
    router.HandleFunc("/api/playerstat", UserStatistic).Methods("GET", "POST")
    router.HandleFunc("/api/motd", GetMotd).Methods("GET")
    router.HandleFunc("/api/gameapi_endpoint", GetGameAPIUri).Methods("GET")

    // Middleware for token checking
    router.Use(TokenChecker)

    // Setup http handler options and start client API server
    if !c.DecodeMetadata.IsDefined("client_api", "use_tls") {
        clientApiLogger.Fatalln("use_tls options is not defined, aborting to start client api server")
    }

    if !checkTLSMandatoryOptions() {
        clientApiLogger.Fatalln("Needed config options for TLS not defined, aborting to start client api server")
    }

    listenaddr, listenport := u.GetListenInformationClientAPI()

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
