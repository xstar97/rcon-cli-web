package config

import (
    "flag"
)

// Constants for routes
var ROUTES = struct {
    RCON         string
    RCON_SERVERS string
    RCON_VERSION string
    LOGS         string
    SAVED        string
}{
    RCON:         "/rcon",
    RCON_SERVERS: "/rcon/servers",
    RCON_VERSION: "/rcon/version",
    LOGS:         "/logs/",
    SAVED:        "/saved",
}

// Constants for commands
var COMMANDS = struct {
    VERSION string
    ENV     string
    CONFIG  string
}{
    VERSION: "--version",
    ENV:     "--env",
    CONFIG:  "--config",
}

// Configuration constants
var CONFIG = struct {
    // Web port
    PORT string
    // Dark/light mode
    MODE string
    // Root path to rcon file
    CLI_ROOT string
    // Root path to rcon.yaml
    CLI_CONFIG string
    // Default rcon env
    CLI_DEFAULT_SERVER string
    // Database type: json
    DB_TYPE string
    DB_JSON_FILE string
}{}

// Parse flags
func init() {
    flag.StringVar(&CONFIG.PORT, "port", "3000", "Web port")
    flag.StringVar(&CONFIG.MODE, "mode", "dark", "Dark/light mode")
    flag.StringVar(&CONFIG.CLI_ROOT, "cli-root", "/app/rcon/rcon", "Root path to rcon file")
    flag.StringVar(&CONFIG.CLI_CONFIG, "cli-config", "/config/rcon.yaml", "Root path to rcon.yaml")
    flag.StringVar(&CONFIG.CLI_DEFAULT_SERVER, "cli-def-server", "default", "Default rcon env")
    flag.StringVar(&CONFIG.DB_TYPE, "db-type", "json", "Database type: json")
    flag.StringVar(&CONFIG.DB_JSON_FILE, "db-json-file", "/config/saved.json", "DB JSON file")
    flag.Parse()
}
