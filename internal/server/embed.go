package server

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var distFS embed.FS

// WebFS returns the embedded web frontend filesystem rooted at "dist".
// Returns nil if the dist directory is empty (dev mode).
func WebFS() fs.FS {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		return nil
	}
	// Check if there's at least an index.html
	if _, err := fs.Stat(sub, "index.html"); err != nil {
		return nil
	}
	return sub
}
