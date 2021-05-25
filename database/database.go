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
    if !c.DecodeMetadata.IsDefined("database", "driver") {
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

    if err != nil {
        dbLogger.Fatalln(err)
    }

    dbLogger.Println("Database initialized")

    err = DB.AutoMigrate(&User{}, &GamesHistory{}, &PlayerList{}, &DHTHosts{})

    if err != nil {
        dbLogger.Fatalln("Error when creating schema: %s\n", err.Error())
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
        "sqlite": []string{
            "sqlite_dbpath",
        },
    }

    // Check mandatory options defined in document
    if options, ok := databaseMandatoryOptions[driverName]; ok {
        for _, v := range options {
            if !c.DecodeMetadata.IsDefined("database", v) {
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

        if c.DecodeMetadata.IsDefined("database", "dbport") {
            databasePort = ":" + strconv.FormatUint(uint64(c.Conf.DB.DatabasePort), 10)
        }

        if c.DecodeMetadata.IsDefined("database", "instance") {
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
