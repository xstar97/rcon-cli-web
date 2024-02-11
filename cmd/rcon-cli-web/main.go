package main

import (
    "fmt"
    "net/http"
    "rcon-cli-web/internal/routes"
    "rcon-cli-web/internal/config"
)

func main() {
    // Serve files from the "public" directory
    fs := http.FileServer(http.Dir("./public"))
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
    http.ListenAndServe(port, nil)
}
