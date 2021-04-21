package logger

import (
    "log"
    "os"
    "fmt"
    "errors"
    "io"
)

var (
    // Errors
    NoLoggerError = errors.New("logger: No logger found by name")

    // Base logger, that available at start of program
    BaseLogger *log.Logger

    // Other loggers, that adds and deletes by other modules
    loggers []*log.Logger
)


func init() {
    // Initialize base logger
    BaseLogger = log.New(os.Stdout, "[Base] ", log.Ldate | log.Ltime | log.Lmicroseconds | log.Lmsgprefix)
}

// Adds a new logger to logger store
func AddNewLogger(name string, output io.Writer, prefix int) (newLog *log.Logger) {
    newLog = log.New(output, fmt.Sprintf("[%s] ", name), prefix)
    loggers = append(loggers, newLog)
    return
}

// Return a loggers that has been added earlier by name
func GetLogger(name string) (*log.Logger, error) {
    for _, logger := range loggers {
        if logger.Prefix() == fmt.Sprintf("[%s] ", name) {
            return logger, nil
        }
    }

    return nil, NoLoggerError
}

// Remove logger from logger store by name
func RemoveLogger(name string) error {
    for i, logger := range loggers {
        if logger.Prefix() == fmt.Sprintf("[%s] ", name) {
            // Fast method to remove item from slice
            // It swaps needed element to the end of slice and reassign then without last element
            loggers[i] = loggers[len(loggers) - 1]
            loggers = loggers[:len(loggers) - 1]
            return nil
        }
    }

    return NoLoggerError
}
