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

    // Check listenaddr and listenport
    if c.Conf.App.ListenAddr == "" || c.Conf.App.ListenPort == "" {
        mainLogger.Fatalln("Error: listenaddr or listenport is not set")
    }

    // Set router handler and start REST API Server
    http.Handle("/", r.Router)
    mainLogger.Fatalln(http.ListenAndServe(fmt.Sprintf(
        "%s:%s",
        c.Conf.App.ListenAddr,
        c.Conf.App.ListenPort,
    ), nil))
}
