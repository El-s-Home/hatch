package handler

import (
	"net/http"

	"github.com/elfoundation/hatch/internal/web"
)

// HandleInspect serves the SPA shell for GET /e/{endpointID}. The client-side
// app reads the endpoint id from the path and renders the live request list.
// All rendering now lives in the shared @hatch/inspector-ui package, so there
// is a single source of truth for how a request is displayed.
func HandleInspect() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		serveIndex(w)
	}
}

func serveIndex(w http.ResponseWriter) {
	html, err := web.Index()
	if err != nil {
		http.Error(w, "UI not built", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(html)
}
