package routes

import (
	"fmt"
	"net/http"
	"embed"
	"strings"
	"rcon-cli-web/internal/config"
)

//go:embed public
var content embed.FS

var (
	staticDir = config.CONFIG.PUBLIC_DIR
)

func rootPath(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			r.URL.Path = fmt.Sprintf("/%s/", staticDir)
		} else {
			firstPathSegment := strings.Split(r.URL.Path, "/")[1]
			if firstPathSegment != staticDir {
				r.URL.Path = fmt.Sprintf("/%s%s", staticDir, r.URL.Path)
			}
		}
		h.ServeHTTP(w, r)
	})
}

func StaticHandler() http.Handler {
	fs := http.FileServer(http.FS(content))
	return rootPath(fs)
}
