package database

import (
    "database/sql"
    "fmt"
    "os"
    "log"

    c "GoChessgameServer/conf"
    "GoChessgameServer/logger"

    _ "github.com/denisenkom/go-mssqldb"
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
        "sqlserver://%s:%s@%s/%s?connection timeout=%d&database=%s",
        c.Conf.DB.Username,
        c.Conf.DB.Password,
        c.Conf.DB.ServerName,
        c.Conf.DB.InstanceName,
        c.Conf.DB.ConnectionTimeout,
        c.Conf.DB.DBName,
    ))
    if err != nil {
        dbLogger.Fatalln(err)
    }

    // Init query executor - go subroutine, that reads queries from channel and executes
    // they synchronically
    initQueryExecutor()
    dbLogger.Println("Database connection initialized")

    // Execute schema creator, that database schema will be created
    // if some elements a absent
    CreateSchemaIfNotExists()
}

// This function creates mandatory schema
func CreateSchemaIfNotExists() {
    // Database logger
    dbLogger, err := logger.GetLogger("Database")
    if err != nil { panic(err) }

    // Execute schema query
    _, err = QueryExecBlocking(`
    IF NOT EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'Users')
    BEGIN
        CREATE TABLE dbo.Users (
            Login VARCHAR(100) PRIMARY KEY,
            Email VARCHAR(100) NOT NULL,
            PasswordHash VARCHAR(256) NOT NULL,
            PasswordHashSalt VARCHAR(256) NOT NULL,
        )
    END
    IF NOT EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'GamesHistory')
    BEGIN
        CREATE TABLE dbo.GamesHistory (
            ID INT PRIMARY KEY IDENTITY(1,1),
            GameStartTimestamp DATETIME NOT NULL,
            GameEndTimestamp DATETIME NOT NULL,
            IsDraw BIT NOT NULL,
            WinnerLogin VARCHAR(100) FOREIGN KEY REFERENCES dbo.Users(Login),
            PlayerOneLogin VARCHAR(100) NOT NULL FOREIGN KEY REFERENCES dbo.Users(Login),
            PlayerTwoLogin VARCHAR(100) NOT NULL FOREIGN KEY REFERENCES dbo.Users(Login)
        )
    END
    `)
    if err != nil {
        dbLogger.Fatalf("Error when executing schema: %s\n", err.Error())
    }

    dbLogger.Println("Database schema creator executed successfully")
}
