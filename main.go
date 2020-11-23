package main

import (
    "GoChessgameServer/logger"
    c "GoChessgameServer/conf"
    _ "GoChessgameServer/util"
    "net/http"
    "fmt"
    "log"
    "io"
    "os"
)

func main() {
    // Register main logger
    mainLogger := logger.AddNewLogger("Application", os.Stdout, log.LstdFlags | log.Lmsgprefix)

    // Check listenaddr and listenport
    if c.Conf.App.ListenAddr == "" || c.Conf.App.ListenPort == "" {
        mainLogger.Fatalln("Error: listenaddr or listenport is not set")
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        mainLogger.Println(r.URL.Path)
        io.WriteString(w, "Hello, world");
    })

    mainLogger.Fatalln(http.ListenAndServe(fmt.Sprintf(
        "%s:%s",
        c.Conf.App.ListenAddr,
        c.Conf.App.ListenPort,
    ), nil))
}
