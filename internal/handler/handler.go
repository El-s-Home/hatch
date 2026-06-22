package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/elfoundation/hatch/internal/store"
)

type Handler struct{ Repo store.Repository }

func New(repo store.Repository) *Handler { return &Handler{Repo: repo} }

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/{endpoint}", h)
	mux.Handle("/{endpoint}/", h)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	eid := extractID(r.URL.Path)
	ctx := context.Background()
	if _, err := h.Repo.GetEndpoint(ctx, eid); err != nil {
		h.Repo.CreateEndpoint(ctx, eid)
	}
	hdr := map[string]string{}
	for k, v := range r.Header { hdr[k] = strings.Join(v, ", ") }
	hdrJSON, _ := json.Marshal(hdr)
	var body []byte
	if r.Body != nil { body, _ = io.ReadAll(r.Body); r.Body.Close() }
	h.Repo.AppendRequest(ctx, eid, &store.Request{
		Method: r.Method, Path: r.URL.Path, Headers: string(hdrJSON),
		Query: r.URL.RawQuery, Body: body,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

func extractID(p string) string {
	t := strings.TrimLeft(p, "/")
	if i := strings.IndexByte(t, '?'); i >= 0 { t = t[:i] }
	if i := strings.IndexByte(t, '/'); i >= 0 { return t[:i] }
	return t
}
