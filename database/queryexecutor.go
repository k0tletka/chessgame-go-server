package database

import (
    "sync"
    "database/sql"
    _ "github.com/denisenkom/go-mssqldb"
    "GoChessgameServer/logger"
)

// Channel that used to be a connector
// between queries and query executor subroutine
var pool = make(chan *queryResult, 50)

// Type that represents object, that
// stores information about query for waiting results
type queryResult struct {
    query string
    args []interface{}
    waitor *sync.WaitGroup
    result *sql.Rows
    resError error
}

func (q *queryResult) GetResults() (*sql.Rows, error) {
    return (*q).result, (*q).resError
}

func (q *queryResult) Waitor() *sync.WaitGroup {
    return (*q).waitor
}

// Interface that will be used to interact
// with other program
type QueryResult interface {
    GetResults() (*sql.Rows, error)
    Waitor() *sync.WaitGroup
}

func initQueryExecutor() {
    // Query executor subroutine start
    go func() {
        dbLogger, err := logger.GetLogger("Database")
        if err != nil { logger.BaseLogger.Fatalln(err) }

        for queryRes := range pool {
            // Execute query
            rows, err := db.Query(queryRes.query, queryRes.args...)
            queryRes.result = rows
            queryRes.resError = err

            // Log execuring query
            dbLogger.Printf("Executed query: %s\n", queryRes.query)

            // Remove added value in waitor
            queryRes.waitor.Done()
        }
    }()
}

// Function to add a new query to pool on execute
func QueryNonBlocking(query string, args ...interface{}) QueryResult {
    waitor := sync.WaitGroup{}

    // Create new queryResult to return then to caller
    queryRes := queryResult{
        query,
        args,
        &waitor,
        nil,
        nil,
    }

    // Add 1 to wait group and push query to executor
    waitor.Add(1)
    pool <- &queryRes

    return &queryRes
}
