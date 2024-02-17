package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"rcon-cli-web/internal/config"
	"rcon-cli-web/internal/routes"
)

func main() {
	// Use the PORT constant from config
	port := fmt.Sprintf(":%s", config.CONFIG.PORT)
    //routes
	rcon := config.ROUTES.RCON
	rconServers := config.ROUTES.RCON_SERVERS
	rconVersion := config.ROUTES.RCON_VERSION
	rconHealth := config.ROUTES.RCON_HEALTH
	logs := config.ROUTES.LOGS
	saved := config.ROUTES.SAVED

	// Define route handlers
	http.HandleFunc(rcon, routes.HandleRcon)
	http.HandleFunc(rconServers, routes.HandleRconServers)
	http.HandleFunc(rconVersion, routes.HandleRconVersion)
	http.HandleFunc(rconHealth, routes.HandleRconHealth)
	http.HandleFunc(logs, routes.HandleLogs)
	http.HandleFunc(saved, routes.HandleSaved)
	http.Handle("/", routes.StaticHandler())

	// Set up signal handling to capture the reason for exit
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start the server
	go func() {
		fmt.Printf("Server is starting on port %s...\n", port)
		if err := http.ListenAndServe(port, nil); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for a signal
	sig := <-sigCh
	log.Printf("Received signal %v. Shutting down...", sig)
}
