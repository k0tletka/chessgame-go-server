package database

import (
    "errors"
    "database/sql"

    "GoChessgameServer/logger"
)

// Errors
var NoColumnsError = errors.New("No columns in result")

// This function is used to parse results from
// *sql.Rows to []interface{}
func ParseRows(rows *sql.Rows) ([][]interface{}, error) {
    defer rows.Close()
    dbLogger, err := logger.GetLogger("Database")
    if err != nil { panic(err) }

    var res [][]interface{}

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
        res = append(res, vals)

        err = rows.Scan(vals...)
        if err != nil {
            dbLogger.Printf("Error occured when parsing row: %s\n", err.Error())
            continue
        }
    }

    if rows.Err() != nil {
        return nil, rows.Err()
    }
    return res, nil
}
