package main

import (
    "fmt"
    "log"
    "os"
    "net/http"

    "GoChessgameServer/logger"
    c "GoChessgameServer/conf"
    _ "GoChessgameServer/store"
    r "GoChessgameServer/router"
    _ "GoChessgameServer/database"
    _ "GoChessgameServer/game"
)

func main() {
    // Register main logger
    mainLogger := logger.AddNewLogger("Application", os.Stdout, log.LstdFlags | log.Lmsgprefix)

    // Check configuration listenport and listenaddr, set default if parameter omited
    var listenaddr string = "0.0.0.0"
    var listenport uint16 = 80

    if c.DecodeMetadata.IsDefined("application", "listenaddr") {
        listenaddr = c.Conf.App.ListenAddr
    }

    if c.DecodeMetadata.IsDefined("application", "listenport") {
        listenport = c.Conf.App.ListenPort
    }

    // Set router handler and start REST API Server
    http.Handle("/", r.Router)
    mainLogger.Fatalln(http.ListenAndServe(fmt.Sprintf(
        "%s:%d",
        listenaddr,
        listenport,
    ), nil))
}
