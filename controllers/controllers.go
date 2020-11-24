package controllers

import (
    "log"
    "os"

    "GoChessgameServer/logger"
)

var contrLogger *log.Logger

func init() {
    // Create logger for controllers module
    contrLogger = logger.AddNewLogger("Controllers", os.Stdout, log.LstdFlags | log.Lmsgprefix)
}
