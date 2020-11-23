package database

import (
    "database/sql"
    _ "github.com/denisenkom/go-mssqldb"
    c "GoChessgameServer/conf"
    "GoChessgameServer/logger"
    "fmt"
    "os"
    "log"
)

// Database variable
var db *sql.DB

func init() {

    // Initialize database logger
    dbLogger := logger.AddNewLogger("Database", os.Stdout, log.LstdFlags | log.Lmsgprefix)

    // Check mandatory credentials and connect to database
    if c.Conf.DB.DBName == "" || c.Conf.DB.Username == "" || c.Conf.DB.Password == "" || c.Conf.DB.ServerName == "" {
        dbLogger.Fatalln("Error: dbname, user or pass is not defined in configuration")
    }

    var err error
    db, err = sql.Open("mssql", fmt.Sprintf(
        "sqlserver://%s:%s@%s/%s?connection timeout=%d",
        c.Conf.DB.Username,
        c.Conf.DB.Password,
        c.Conf.DB.ServerName,
        c.Conf.DB.InstanceName,
        c.Conf.DB.ConnectionTimeout,
    ))
    if err != nil {
        dbLogger.Fatalln(err)
    }

    // Init query executor - go subroutine, that reads queries from channel and executes
    // they synchronically
    initQueryExecutor()
    dbLogger.Println("Database initialized")
}
