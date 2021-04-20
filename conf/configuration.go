package conf

import (
    "github.com/BurntSushi/toml"
    "GoChessgameServer/logger"
    "log"
    "os"
)

type Configuration struct {
    DB Database `toml:"database"`
    App Application `toml:"application"`
}

type Database struct {
    ServerName string `toml:"server"`
    InstanceName string `toml:"instance"`
    DBName string `toml:"dbname"`
    Username string `toml:"user"`
    Password string `toml:"pass"`
    ConnectionTimeout int `toml:"connection_timeout"`
}

type Application struct {
    ListenAddr string `toml:"listenaddr"`
    ListenPort string `toml:"listenport"`
    MarkdownFile string `toml:"markdown_file"`
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

    // Decode configuration
    if _, err := toml.DecodeFile(conffile, &Conf); err != nil {
        confLogger.Fatalln(err)
    }

    confLogger.Printf("Configuration from %s loaded\n", conffile)
}
