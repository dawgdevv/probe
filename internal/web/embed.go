package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:static
var embeddedFiles embed.FS

// StaticFS returns the embedded static filesystem for serving the UI
func StaticFS() (http.FileSystem, error) {
	sub, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		return nil, err
	}
	return http.FS(sub), nil
}

// IndexHTML returns the contents of index.html for SPA fallback
func IndexHTML() ([]byte, error) {
	return embeddedFiles.ReadFile("static/index.html")
}
