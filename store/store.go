package store

import (
    "log"
    "os"
    "io/ioutil"

    "GoChessgameServer/logger"
    c "GoChessgameServer/conf"
)

var (
    // String variable for markdown motd storing
    MotdString string

    // Utility
    storeLogger *log.Logger
)

func init() {
    // Create store logger
    storeLogger = logger.AddNewLogger("Store", os.Stdout, log.LstdFlags | log.Lmsgprefix)

    // Load markdown file
    mdFile := c.Conf.CAPI.MarkdownFile

    if mdFile == "" {
        storeLogger.Println("Markdown file is not set, skipping...")
        MotdString = "Motd file is not set on the server."
    } else {
        mdFd, err := os.OpenFile(mdFile, os.O_RDONLY | os.O_CREATE, 0755)

        if err != nil {
            storeLogger.Fatalln(err)
        }
        defer mdFd.Close()

        // Read markdown into string variable
        readedMotd, err := ioutil.ReadAll(mdFd)

        if err != nil {
            storeLogger.Fatalln(err)
        }

        MotdString = string(readedMotd)
    }
}
