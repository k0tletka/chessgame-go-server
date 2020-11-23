package main

import (
    "net/http"
    "fmt"
    "log"
    "os"

    "GoChessgameServer/logger"
    c "GoChessgameServer/conf"
    _ "GoChessgameServer/util"
    _ "GoChessgameServer/database"

    _ "github.com/gorilla/mux"
)

func main() {
    // Register main logger
    mainLogger := logger.AddNewLogger("Application", os.Stdout, log.LstdFlags | log.Lmsgprefix)

    // Check listenaddr and listenport
    if c.Conf.App.ListenAddr == "" || c.Conf.App.ListenPort == "" {
        mainLogger.Fatalln("Error: listenaddr or listenport is not set")
    }

    // Define router, bind handlers on controllers and define middleware
    // for authetication handling

    mainLogger.Fatalln(http.ListenAndServe(fmt.Sprintf(
        "%s:%s",
        c.Conf.App.ListenAddr,
        c.Conf.App.ListenPort,
    ), nil))
}
