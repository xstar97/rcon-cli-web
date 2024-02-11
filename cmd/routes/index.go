// routes/index.go

package routes

import (
    "log"
    "net/http"
    "os"
    "rcon-cli-web/config"
)

func MainIndexRoute() {
    // Check if the public directory exists
    publicDir := config.CONFIG.PUBLIC_DIR
    if _, err := os.Stat(publicDir); os.IsNotExist(err) {
        log.Fatalf("The %s directory does not exist", publicDir)
    }

    // Serve files from the public directory
    fs := http.FileServer(http.Dir(publicDir))
    http.Handle("/", fs)
}