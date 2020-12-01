package longpollmanagers

import (
    "log"
    "os"

    "GoChessgameServer/logger"
)

// Longpoll managers logger
var lmLogger *log.Logger

func init() {
    // Create new longpoll manager logger
    lmLogger = logger.AddNewLogger("LongpollManager", os.Stdout, log.LstdFlags | log.Lmsgprefix)
}
