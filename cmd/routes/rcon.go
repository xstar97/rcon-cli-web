package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"rcon-cli-web/config"
	"strings"
	"gopkg.in/yaml.v2"
    "log"
	"os/exec"
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

func HandleRcon(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        log.Println("Method not allowed")
        return
    }

    // Decode the request body into RconRequest struct
    var req RconRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Failed to decode request body", http.StatusBadRequest)
        log.Println("Failed to decode request body:", err)
        return
    }

    log.Println("Received Rcon request:", req)

    // Execute the shell command
    output, err := ExecuteShellCommand(req.Server, req.Command)
    if err != nil {
        http.Error(w, "Failed to run rcon-cli: "+err.Error(), http.StatusInternalServerError)
        log.Println("Failed to run rcon-cli:", err)
        return
    }

    log.Println("Rcon command executed successfully")

    // Prepare the response object
    resp := RconResponse{
        Server:  req.Server,
        Command: req.Command,
        Output:  strings.TrimSpace(string(output)), // Trim leading/trailing whitespace
    }

    // Log the response object
    log.Println("Rcon response prepared:", resp)

    // Set the content type to application/json
    w.Header().Set("Content-Type", "application/json")

    // Encode the response to JSON and write it to the response writer
    if err := json.NewEncoder(w).Encode(resp); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        log.Println("Failed to encode response:", err)
        return
    }

    log.Println("Rcon response sent successfully")
}

func HandleRconServers(w http.ResponseWriter, r *http.Request) {
    // Read server configurations from the file
    configFile := config.CONFIG.CLI_CONFIG
    configData, err := ioutil.ReadFile(configFile)
    if err != nil {
        http.Error(w, "Failed to read server configurations", http.StatusInternalServerError)
        log.Println("Failed to read server configurations:", err)
        return
    }

    log.Println("Server configurations read successfully")

    // Parse the YAML configuration data
    var servers map[string]map[string]string
    err = yaml.Unmarshal(configData, &servers)
    if err != nil {
        http.Error(w, "Failed to parse server configurations", http.StatusInternalServerError)
        log.Println("Failed to parse server configurations:", err)
        return
    }

    log.Println("Server configurations parsed successfully")

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

    log.Println("Server data extracted successfully")

    // Set the content type to application/json
    w.Header().Set("Content-Type", "application/json")

    // Encode the server data to JSON and write it to the response writer
    if err := json.NewEncoder(w).Encode(serverData); err != nil {
        http.Error(w, "Failed to encode server data", http.StatusInternalServerError)
        log.Println("Failed to encode server data:", err)
        return
    }

    log.Println("Server data encoded and sent successfully")
}

func HandleRconVersion(w http.ResponseWriter, r *http.Request) {
    // Get the latest release information from GitHub API
    resp, err := http.Get("https://api.github.com/repos/gorcon/rcon-cli/releases/latest")
    if err != nil {
        http.Error(w, "Failed to fetch latest version information", http.StatusInternalServerError)
        log.Println("Failed to fetch latest version information:", err)
        return
    }
    defer resp.Body.Close()

    log.Println("Latest release information fetched successfully")

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, "Failed to read latest version information", http.StatusInternalServerError)
        log.Println("Failed to read latest version information:", err)
        return
    }

    log.Println("Latest release information read successfully")

    // Parse the JSON response
    var releaseInfo struct {
        TagName string `json:"tag_name"`
    }
    err = json.Unmarshal(body, &releaseInfo)
    if err != nil {
        http.Error(w, "Failed to parse latest version information", http.StatusInternalServerError)
        log.Println("Failed to parse latest version information:", err)
        return
    }

    log.Println("Latest release information parsed successfully")

    // Extract the latest version
    latestVersion := strings.TrimPrefix(releaseInfo.TagName, "v")

    // Get the current version
    output, err := ExecuteShellCommand(config.CONFIG.CLI_DEFAULT_SERVER, config.COMMANDS.VERSION)
    if err != nil {
        http.Error(w, "Failed to run rcon-cli: "+err.Error(), http.StatusInternalServerError)
        log.Println("Failed to run rcon-cli:", err)
        return
    }

    log.Println("Current version retrieved successfully")

    // Extract the current version from the output
    currentVersion := strings.TrimSpace(strings.TrimPrefix(string(output), "rcon version"))

    // Determine if an update is available
    updateAvailable := latestVersion != currentVersion

    log.Println("Update availability checked")

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
        log.Println("Failed to encode response:", err)
        return
    }

    log.Println("Response encoded and sent successfully")
}
