package config

import (
    "flag"
    "log"
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
    // Public directory
    PUBLIC_DIR string
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
    flag.StringVar(&CONFIG.PUBLIC_DIR, "public-dir", "/app/public", "Public directory")
    flag.Parse()

    // Log the set flags
    log.Printf("Web port: %s", CONFIG.PORT)
    log.Printf("Dark/light mode: %s", CONFIG.MODE)
    log.Printf("Root path to rcon file: %s", CONFIG.CLI_ROOT)
    log.Printf("Root path to rcon.yaml: %s", CONFIG.CLI_CONFIG)
    log.Printf("Default rcon env: %s", CONFIG.CLI_DEFAULT_SERVER)
    log.Printf("Database type: %s", CONFIG.DB_TYPE)
    log.Printf("DB JSON file: %s", CONFIG.DB_JSON_FILE)
    log.Printf("Public directory: %s", CONFIG.PUBLIC_DIR)
}
