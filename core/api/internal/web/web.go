// Package web embeds the built Vite SPA so the whole product ships as a
// single Go binary — no Node runtime required to self-host.
package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:dist
var dist embed.FS

// distFS is the dist/ subtree rooted so paths look like "index.html",
// "assets/...". If the SPA has not been built, sub() still succeeds with an
// effectively empty tree and the handlers below degrade gracefully.
func distFS() fs.FS {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		return dist
	}
	return sub
}

// AssetsHandler serves hashed static assets under /assets/.
func AssetsHandler() http.Handler {
	return http.FileServer(http.FS(distFS()))
}

// Index returns the SPA entry document.
func Index() ([]byte, error) {
	return fs.ReadFile(distFS(), "index.html")
}
