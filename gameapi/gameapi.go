package gameapi

import (
    "log"
    "os"

    "GoChessgameServer/logger"
)

var (
    // Logger for Game API module
    gameApiLogger *log.Logger
)

func init() {
    // Create logger for controllers module
    gameApiLogger = logger.AddNewLogger("GameAPI", os.Stdout, log.LstdFlags | log.Lmsgprefix)
}
