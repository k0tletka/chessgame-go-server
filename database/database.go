package database

import (
    "fmt"
    "os"
    "log"
    "strconv"

    c "GoChessgameServer/conf"
    "GoChessgameServer/logger"

    "gorm.io/gorm"

    // Drivers for gorn
    sqlserverdriver "gorm.io/driver/sqlserver"
    sqlitedriver "gorm.io/driver/sqlite"
)

var (
    // Database variable
    DB *gorm.DB

    // Logger
    dbLogger *log.Logger
)

func init() {

    // Initialize database logger
    dbLogger = logger.AddNewLogger("Database", os.Stdout, log.LstdFlags | log.Lmsgprefix)

    // Check mandatory credentials and connect to database
    if !c.DecodeMetadata.IsDefined("database.driver") {
        dbLogger.Fatalln("Database driver is not defined in config file, aborting...")
    }

    if !checkMandatoryOptions(c.Conf.DB.DatabaseDriver) {
        dbLogger.Fatalln("One or multiple mandatory options in config not defined, aborting...")
    }

    // Get gorm dialector, based on driver
    dbDialector := getAppropriateDialector(c.Conf.DB.DatabaseDriver)

    // Initialize connection
    var err error
    DB, err = gorm.Open(dbDialector, &gorm.Config{})

    dbLogger.Println("Database initialized")

    // Connect to the database with
    // Init query executor - go subroutine, that reads queries from channel and executes
    // they synchronically
    //initQueryExecutor()

    // Execute schema creator, that database schema will be created
    // if some elements a absent
    CreateSchemaIfNotExists()
}

// This function creates mandatory schema
func CreateSchemaIfNotExists() {

    // Execute schema query
    /*_, err = QueryExecBlocking(`
    IF NOT EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'Users')
    BEGIN
        CREATE TABLE dbo.Users (
            Login VARCHAR(100) PRIMARY KEY,
            Email VARCHAR(100) NOT NULL,
            IsAdmin BIT NOT NULL DEFAULT 0,
            PasswordHash VARBINARY(1000) NOT NULL,
            PasswordHashSalt VARBINARY(1000) NOT NULL,
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
    */

    if !DB.Migrator().HasTable(&User{}) {
        DB.Migrator().CreateTable(&User{})
    }

    if !DB.Migrator().HasTable(&GamesHistory{}) {
        DB.Migrator().CreateTable(&GamesHistory{})
    }

    dbLogger.Println("Database schema creator executed successfully")
}

// Check mandatory options for different drivers
func checkMandatoryOptions(driverName string) bool {

    databaseMandatoryOptions := map[string][]string{
        "sqlserver": []string{
            "server",
            "dbname",
            "user",
            "pass",
        },
        "sqlite3": []string{
            "sqlite_dbpath",
        },
    }

    // Check mandatory options defined in document
    if options, ok := databaseMandatoryOptions[driverName]; ok {
        for _, v := range options {
            if !c.DecodeMetadata.IsDefined("database." + v) {
                return false
            }
        }

        return true
    }

    return false
}

// Returns dialectors based on passed driver name.
// If invalid driver name given, func returns nil
func getAppropriateDialector(driverName string) gorm.Dialector {

    switch driverName {
    case "sqlserver":
        var databasePort string
        var instanceName string

        if c.DecodeMetadata.IsDefined("database.dbport") {
            databasePort = ":" + strconv.FormatUint(uint64(c.Conf.DB.DatabasePort), 10)
        }

        if c.DecodeMetadata.IsDefined("database.instance") {
            instanceName = "/" + c.Conf.DB.InstanceName
        }

        return sqlserverdriver.Open(fmt.Sprintf(
            "sqlserver://%s:%s@%s%s%s?database=%s",
            c.Conf.DB.Username,
            c.Conf.DB.Password,
            c.Conf.DB.ServerName,
            databasePort,
            instanceName,
            c.Conf.DB.DBName,
        ))
    case "sqlite":
        return sqlitedriver.Open(c.Conf.DB.SqliteDatabasePath)
    default:
        return nil
    }

}
