package websocket

import (
    "log"
    "os"

    "GoChessgameServer/logger"
)

var (
    // Logger for websocket module
    websocketLogger *log.Logger
)

func init() {
    // Create new logger
    websocketLogger = logger.AddNewLogger("Websocket", os.Stdout, log.LstdFlags | log.Lmsgprefix)
}
