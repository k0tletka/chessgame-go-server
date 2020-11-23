package conf

import (
    "github.com/BurntSushi/toml"
    "GoChessgameServer/logger"
    "log"
    "os"
    "bytes"
    "io"
)

type Configuration struct {
    DB Database `toml:"database"`
    App Application `toml:"application"`
}

type Database struct {
    DBName string `toml:"dbname"`
    Username string `toml:"user"`
    Password string `toml:"pass"`
}

type Application struct {
    ListenAddr string `toml:"listenaddr"`
    ListenPort string `toml:"listenport"`
}

// Conf variable
var Conf Configuration

func init() {

    // Initialize configuration
    confLogger := logger.AddNewLogger("Configuration", os.Stdout, log.LstdFlags | log.Lmsgprefix)
    conffile := os.Getenv("CONFLOCATION")
    if conffile == "" {
        conffile = "configuration.toml"
    }

    // Open configuration file
    readStream, err := os.OpenFile(conffile, os.O_RDONLY | os.O_CREATE, 0755)
    if err != nil {
        confLogger.Fatalln(err)
    }

    // Read data
    buffer := bytes.Buffer{}
    _, err = io.Copy(&buffer, readStream)
    if err != nil {
        confLogger.Fatalln(err)
    }

    // Decode configuration
    if _, err = toml.Decode(buffer.String(), &Conf); err != nil {
        confLogger.Fatalln(err)
    }

    confLogger.Printf("Configuration from %s loaded\n", conffile)
}
