package main

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	webui "github.com/VictoriaMetrics/mcp-victoriametrics/web"
)

func spaHandler() http.Handler {
	dist, err := fs.Sub(webui.DistFS, "dist")
	if err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "UI build not available", http.StatusInternalServerError)
		})
	}

	fileServer := http.FileServer(http.FS(dist))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		cleanPath := path.Clean(r.URL.Path)
		if cleanPath == "/" {
			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)
			return
		}

		if strings.Contains(path.Base(cleanPath), ".") {
			if _, err := fs.Stat(dist, strings.TrimPrefix(cleanPath, "/")); err == nil {
				r.URL.Path = cleanPath
				fileServer.ServeHTTP(w, r)
				return
			}
			http.NotFound(w, r)
			return
		}

		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})
}
