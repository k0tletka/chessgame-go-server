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
    DatabaseDriver      string  `toml:"driver"`
    ServerName          string  `toml:"server"`
    DatabasePort        uint16  `toml:"dbport"`
    InstanceName        string  `toml:"instance"`
    DBName              string  `toml:"dbname"`
    SqliteDatabasePath  string  `toml:"sqlite_dbpath"`
    Username            string  `toml:"user"`
    Password            string  `toml:"pass"`
    ConnectionTimeout   int     `toml:"connection_timeout"`
}

type Application struct {
    ListenAddr          string  `toml:"listenaddr"`
    ListenPort          uint16    `toml:"listenport"`
    MarkdownFile        string  `toml:"markdown_file"`
}

var (
    // Variable for configuration store
    Conf Configuration

    // TOML decoding metadata
    DecodeMetadata toml.MetaData
)

func init() {

    // Initialize configuration
    confLogger := logger.AddNewLogger("Configuration", os.Stdout, log.LstdFlags | log.Lmsgprefix)

    conffile, defined := os.LookupEnv("CONFLOCATION")
    if !defined {
        conffile = "configuration.toml"
    }

    // Decode configuration
    var err error

    if DecodeMetadata, err = toml.DecodeFile(conffile, &Conf); err != nil {
        confLogger.Fatalln(err)
    }

    confLogger.Printf("Configuration from %s loaded\n", conffile)
}
