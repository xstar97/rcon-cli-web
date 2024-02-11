package routes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"rcon-cli-web/config"
	"strings"
	"gopkg.in/yaml.v2"
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

	// Read server configurations from the file
	configFile := config.CONFIG.CLI_CONFIG
	fmt.Println("Config file:", configFile)

	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		http.Error(w, "Failed to read server configurations", http.StatusInternalServerError)
		fmt.Println("Error reading config file:", err)
		return
	}

	fmt.Println("Server configurations read successfully")

	// Parse the YAML configuration data
	var servers map[string]map[string]string
	err = yaml.Unmarshal(configData, &servers)
	if err != nil {
		http.Error(w, "Failed to parse server configurations", http.StatusInternalServerError)
		fmt.Println("Error parsing config data:", err)
		return
	}

	fmt.Println("Server configurations parsed successfully")

	// Get the log file path for the specified server
	logFile, ok := servers[server]["log"]
	if !ok {
		http.Error(w, "Log file not found for server "+server, http.StatusNotFound)
		fmt.Println("Log file not found for server:", server)
		return
	}

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
