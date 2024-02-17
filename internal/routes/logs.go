package routes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"rcon-cli-web/internal/config"
	"strings"
)

func HandleLogs(w http.ResponseWriter, r *http.Request) {
	// Extract server name from URL path
	server := strings.TrimPrefix(r.URL.Path, "/logs/")
	fmt.Println("Server:", server)

	// If server name is empty, default to the CLI default server
	if server == "" {
		server = config.CONFIG.CLI_DEFAULT_SERVER
		fmt.Println("Using default server:", server)
	}

	// Read server configuration using GetServer
	serverConfig, err := config.GetServer(server)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get server '%s' configuration", server), http.StatusInternalServerError)
		fmt.Println("Error getting server configuration:", err)
		return
	}

	fmt.Println("Server configurations read successfully")

	// Get the log file path for the specified server
	logFile := serverConfig.Log
	fmt.Println("Log file:", logFile)

	// Read the contents of the log file
	logData, err := ioutil.ReadFile(logFile)
	if err != nil {
		http.Error(w, "Failed to read log file", http.StatusInternalServerError)
		fmt.Println("Error reading log file:", err)
		return
	}

	fmt.Println("Log file read successfully")

	// Set the content type to text/plain
	w.Header().Set("Content-Type", "text/plain")

	// Write the log data to the response writer
	w.Write(logData)

	fmt.Println("Log data sent successfully")
}
