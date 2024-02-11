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
    // Check if the public directory exists
    if _, err := os.Stat(config.CONFIG.PUBLIC_DIR); os.IsNotExist(err) {
        log.Fatalf("The %s directory does not exist", config.CONFIG.PUBLIC_DIR)
    }

    // Serve files from the public directory
    fs := http.FileServer(http.Dir(config.CONFIG.PUBLIC_DIR))
    http.Handle("/", fs)

    // Use the PORT constant from config
    port := fmt.Sprintf(":%s", config.CONFIG.PORT)
    rcon := config.ROUTES.RCON
    rconServers := config.ROUTES.RCON_SERVERS
    rconVersion := config.ROUTES.RCON_VERSION
    logs := config.ROUTES.LOGS
    saved := config.ROUTES.SAVED

    // Define route handlers
    http.HandleFunc(rcon, routes.HandleRcon)
    http.HandleFunc(rconServers, routes.HandleRconServers)
    http.HandleFunc(rconVersion, routes.HandleRconVersion)
    http.HandleFunc(logs, routes.HandleLogs)
    http.HandleFunc(saved, routes.HandleSaved)

    fmt.Printf("Server is listening on port %s...\n", port)

    // Set up signal handling to capture the reason for exit
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        sig := <-sigCh
        log.Printf("Received signal %v. Exiting...", sig)
        os.Exit(0)
    }()

    // Start the server
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Printf("Server error: %v", err)
    }
}
