package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"
	"rcon-cli-web/internal/config"
	"strings"
	"gopkg.in/yaml.v2"
)

type RconRequest struct {
    Server  string `json:"server"`
    Command string `json:"command"`
}

type RconResponse struct {
    Server  string `json:"server"`
    Command string `json:"command"`
    Output  string `json:"output"`
}

type RconVersionResponse struct {
	LatestVersion   string `json:"latestVersion"`
	CurrentVersion  string `json:"currentVersion"`
	UpdateAvailable bool   `json:"updateAvailable"`
}

func HandleRcon(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Decode the request body into RconRequest struct
    var req RconRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Failed to decode request body", http.StatusBadRequest)
        return
    }

    // Execute the shell command
    output, err := ExecuteShellCommand(req.Server, req.Command)
    if err != nil {
        http.Error(w, "Failed to run rcon-cli: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Prepare the response object
    resp := RconResponse{
        Server:  req.Server,
        Command: req.Command,
        Output:  strings.TrimSpace(string(output)), // Trim leading/trailing whitespace
    }

    // Set the content type to application/json
    w.Header().Set("Content-Type", "application/json")

    // Encode the response to JSON and write it to the response writer
    if err := json.NewEncoder(w).Encode(resp); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
}


func HandleRconServers(w http.ResponseWriter, r *http.Request) {
    // Read server configurations from the file
    configFile := config.CONFIG.CLI_CONFIG
    configData, err := ioutil.ReadFile(configFile)
    if err != nil {
        http.Error(w, "Failed to read server configurations", http.StatusInternalServerError)
        return
    }

    // Parse the YAML configuration data
    var servers map[string]map[string]string
    err = yaml.Unmarshal(configData, &servers)
    if err != nil {
        http.Error(w, "Failed to parse server configurations", http.StatusInternalServerError)
        return
    }

    // Extract server name and type from each configuration
    var serverData []struct {
        Server string `json:"server"`
        Type   string `json:"type"`
    }
    for name, config := range servers {
        server := struct {
            Server string `json:"server"`
            Type   string `json:"type"`
        }{
            Server: name,
            Type:   config["type"],
        }
        serverData = append(serverData, server)
    }

    // Set the content type to application/json
    w.Header().Set("Content-Type", "application/json")

    // Encode the server data to JSON and write it to the response writer
    if err := json.NewEncoder(w).Encode(serverData); err != nil {
        http.Error(w, "Failed to encode server data", http.StatusInternalServerError)
        return
    }
}

// ExecuteShellCommand executes a shell command and returns its output
func ExecuteShellCommand(server, command string) ([]byte, error) {
    // Set the command to execute
    cmd := exec.Command(config.CONFIG.CLI_ROOT, config.COMMANDS.CONFIG, config.CONFIG.CLI_CONFIG, config.COMMANDS.ENV, server, command)
    
    // Capture the output of the command
    output, err := cmd.CombinedOutput()
    if err != nil {
        return nil, err
    }

    return output, nil
}

func HandleRconVersion(w http.ResponseWriter, r *http.Request) {
	// Get the latest release information from GitHub API
	resp, err := http.Get("https://api.github.com/repos/gorcon/rcon-cli/releases/latest")
	if err != nil {
		http.Error(w, "Failed to fetch latest version information", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read latest version information", http.StatusInternalServerError)
		return
	}

	// Parse the JSON response
	var releaseInfo struct {
		TagName string `json:"tag_name"`
	}
	err = json.Unmarshal(body, &releaseInfo)
	if err != nil {
		http.Error(w, "Failed to parse latest version information", http.StatusInternalServerError)
		return
	}

	// Extract the latest version
	latestVersion := strings.TrimPrefix(releaseInfo.TagName, "v")

	// Get the current version
    output, err := ExecuteShellCommand(config.CONFIG.CLI_DEFAULT_SERVER, config.COMMANDS.VERSION)
    if err != nil {
        http.Error(w, "Failed to run rcon-cli: "+err.Error(), http.StatusInternalServerError)
        return
    }
    // Extract the current version from the output
    currentVersion := strings.TrimSpace(strings.TrimPrefix(string(output), "rcon version"))

	// Determine if an update is available
	updateAvailable := latestVersion != currentVersion

	// Construct the response
	response := RconVersionResponse{
		LatestVersion:   latestVersion,
		CurrentVersion:  currentVersion,
		UpdateAvailable: updateAvailable,
	}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the response to JSON and write it to the response writer
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
