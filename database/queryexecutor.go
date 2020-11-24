package database

import (
    "sync"
    "database/sql"

    //"GoChessgameServer/logger"

    _ "github.com/denisenkom/go-mssqldb"
)

// Channels that used to be a connector
// between queries and query executor subroutines
var poolquery = make(chan *queryResult, 50)
var poolexecquery = make(chan *queryExecResult, 50)

// Mutex that will controll query execution in
// different subroutines
var executorMutex = sync.Mutex{}

// Types that represents object, that
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

type queryExecResult struct {
    query string
    args []interface{}
    waitor *sync.WaitGroup
    result *sql.Result
    resError error
}

func (q *queryExecResult) GetResults() (*sql.Result, error) {
    return (*q).result, (*q).resError
}

func (q *queryExecResult) Waitor() *sync.WaitGroup {
    return (*q).waitor
}

// Interfaces that will be used to interact
// with other program
type QueryResult interface {
    GetResults() (*sql.Rows, error)
    Waitor() *sync.WaitGroup
}

type QueryExecResult interface {
    GetResults() (*sql.Result, error)
    Waitor() *sync.WaitGroup
}

func initQueryExecutor() {
    //dbLogger, err := logger.GetLogger("Database")
    //if err != nil { logger.BaseLogger.Fatalln(err) }

    // Query executor subroutine start
    go func() {
        for queryRes := range poolquery {
            // Execute query
            executorMutex.Lock()
            rows, err := db.Query(queryRes.query, queryRes.args...)
            executorMutex.Unlock()

            queryRes.result = rows
            queryRes.resError = err

            // Log executing query (just for debugging)
            //dbLogger.Printf("Executed query: %s\n", queryRes.query)

            // Remove added value in waitor
            queryRes.waitor.Done()
        }
    }()

    // Query exec executor subrouting start
    go func() {
        for queryExecRes := range poolexecquery {
            // Execute query
            executorMutex.Lock()
            res, err := db.Exec(queryExecRes.query, queryExecRes.args...)
            executorMutex.Unlock()

            queryExecRes.result = &res
            queryExecRes.resError = err

            // Log executing query (just for debugging)
            //dbLogger.Printf("Executed query: %s\n", queryExecRes.query)

            // Remove added value in waitor
            queryExecRes.waitor.Done()
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
    poolquery <- &queryRes

    return &queryRes
}

// Function to add a new exec query to pool on execute
func QueryExecNonBlocking(query string, args ...interface{}) QueryExecResult {
    waitor := sync.WaitGroup{}

    // Create new queryExecResult to return then to caller
    queryExecRes := queryExecResult{
        query,
        args,
        &waitor,
        nil,
        nil,
    }

    // Add 1 to wait group and push exec query to executor
    waitor.Add(1)
    poolexecquery <- &queryExecRes

    return &queryExecRes
}
