package database

import (
    "errors"
    "database/sql"

    "GoChessgameServer/logger"
)

// Type of rows results
type RowsResult []map[string]interface{}

// Errors
var NoColumnsError = errors.New("No columns in result")

// This function is used to parse results from
// *sql.Rows to []interface{}
func ParseRows(rows *sql.Rows) (*RowsResult, error) {
    defer rows.Close()
    dbLogger, err := logger.GetLogger("Database")
    if err != nil { panic(err) }

    var res RowsResult

    // Get rows columns
    cols, err := rows.Columns()
    if err != nil {
        return nil, err
    }
    if cols == nil {
        return nil, NoColumnsError
    }

    for rows.Next() {
        vals := make([]interface{}, len(cols))
        err = rows.Scan(vals...)
        if err != nil {
            dbLogger.Printf("Error occured when parsing row: %s\n", err.Error())
            continue
        }

        mapvals := make(map[string]interface{})
        res = append(res, mapvals)

        for i, colName := range cols {
            mapvals[colName] = vals[i]
        }
    }

    if rows.Err() != nil {
        return nil, rows.Err()
    }
    return &res, nil
}

// This function executes queries with waiting
// results
func QueryBlocking(query string, args ...interface{}) (*RowsResult, error) {
    // Execute query
    queryRes := QueryNonBlocking(query, args...)
    queryRes.Waitor().Wait()

    rows, err := queryRes.GetResults()
    if err != nil {
        return nil, err
    }

    // Parse rows
    rowsResults, err := ParseRows(rows)
    if err != nil {
        return nil, err
    }

    return rowsResults, nil
}

// This function executes exec queries with waiting
// results
func QueryExecBlocking(query string, args ...interface{}) (*sql.Result, error) {
    // Execute exec query
    queryRes := QueryExecNonBlocking(query, args...)
    queryRes.Waitor().Wait()

    result, err := queryRes.GetResults()
    if err != nil {
        return nil, err
    }

    return result, nil
}
