package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    "sync"
    "net/http"
    "context"

    // Trigger init methods for packages
    "GoChessgameServer/logger"
    _ "GoChessgameServer/conf"
    _ "GoChessgameServer/store"
    _ "GoChessgameServer/database"

    // Servers
    clientAPI "GoChessgameServer/clientapi"
    gameAPI "GoChessgameServer/gameapi"
    DHTAPI "GoChessgameServer/dht"
)

var (
    // List of running servers
    servers = []*http.Server{}
)

func main() {

    // Defined logger for main application
    mainLogger := logger.AddNewLogger("Application", os.Stdout, log.LstdFlags | log.Lmsgprefix)

    // Define handling of process signals
    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

    // Define channel for server-object grabbing from servers goroutines
    serverChan := make(chan *http.Server, 1)

    go func() {
        for server := range serverChan {
            servers = append(servers, server)
        }
    }()

    // sync.WaitGroup for waiting to start all servers
    srvWaitor := &sync.WaitGroup{}
    srvWaitor.Add(3) // Client API, Game API

    // Start servers
    go clientAPI.InitializeClientAPIServer(srvWaitor, serverChan)

    go gameAPI.InitializeGameAPIServer(srvWaitor, serverChan)

    go DHTAPI.InitializeDHT(srvWaitor, serverChan)

    // Wait for servers to start
    srvWaitor.Wait()
    close(serverChan)

    mainLogger.Println("All servers has started successfully")

    // Signal handling
    for inputSignal := range signalChan {
        switch inputSignal {
        case syscall.SIGHUP:
            // Ignore signal
            continue
        case syscall.SIGINT, syscall.SIGTERM:
            mainLogger.Println("Terminating...")

            // Close all servers gracefully
            for _, server := range servers {
                if err := server.Shutdown(context.Background()); err != nil {
                    mainLogger.Fatalln("Error when stoping server: ", err.Error())
                }
            }

            // Exit application
            os.Exit(0)
        }
    }
}
