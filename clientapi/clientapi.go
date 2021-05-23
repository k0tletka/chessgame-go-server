package clientapi

import (
    "log"
    "os"

    "GoChessgameServer/logger"
)

var (
    clientApiLogger *log.Logger
)

func init() {
    // Create new logger for clientapi package
    clientApiLogger = logger.AddNewLogger("ClientAPI", os.Stdout, log.LstdFlags | log.Lmsgprefix)
}
