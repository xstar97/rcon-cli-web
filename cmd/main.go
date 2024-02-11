package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "rcon-cli-web/config"
    "rcon-cli-web/routes"
)

func main() {
    // Use the PORT constant from config
    port := fmt.Sprintf(":%s", config.CONFIG.PORT)
    rcon := config.ROUTES.RCON
    rconServers := config.ROUTES.RCON_SERVERS
    rconVersion := config.ROUTES.RCON_VERSION
    logs := config.ROUTES.LOGS
    saved := config.ROUTES.SAVED

    // Initialize the index route for serving static files
    routes.MainIndexRoute()

    // Define route handlers
    http.HandleFunc(rcon, routes.HandleRcon)
    http.HandleFunc(rconServers, routes.HandleRconServers)
    http.HandleFunc(rconVersion, routes.HandleRconVersion)
    http.HandleFunc(logs, routes.HandleLogs)
    http.HandleFunc(saved, routes.HandleSaved)

    // Start the server
    go func() {
        fmt.Printf("Server is starting on port %s...\n", port)
        if err := http.ListenAndServe(port, nil); err != nil {
            log.Fatalf("Server failed to start: %v", err)
        }
    }()

    // Set up signal handling to capture the reason for exit
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

    // Wait for a signal
    sig := <-sigCh
    log.Printf("Received signal %v. Shutting down...", sig)
}
