package config

import (
    "os"
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
}{
    PORT:               getEnv("PORT", "3000"),
    MODE:               getEnv("MODE", "dark"),
    CLI_ROOT:           getEnv("CLI_ROOT", "/code/examples/rcon/rcon"),
    CLI_CONFIG:         getEnv("CLI_CONFIG", "/code/examples/rcon/rcon.yaml"),
    CLI_DEFAULT_SERVER: getEnv("CLI_DEFAULT_SERVER", "default"),
    DB_TYPE:            getEnv("DB_TYPE", "json"),
    DB_JSON_FILE:       getEnv("DB_JSON_FILE", "saved.json"),
}

// Get environment variable value or default
func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}
