package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/elfoundation/hatch/internal/store"
)

// sseHub manages SSE subscribers per endpoint.
type sseHub struct {
	mu   sync.RWMutex
	subs map[string]map[chan []byte]struct{}
}

var hub = &sseHub{
	subs: make(map[string]map[chan []byte]struct{}),
}

func (h *sseHub) subscribe(endpointID string) chan []byte {
	ch := make(chan []byte, 64)
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.subs[endpointID] == nil {
		h.subs[endpointID] = make(map[chan []byte]struct{})
	}
	h.subs[endpointID][ch] = struct{}{}
	return ch
}

func (h *sseHub) unsubscribe(endpointID string, ch chan []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.subs[endpointID], ch)
	if len(h.subs[endpointID]) == 0 {
		delete(h.subs, endpointID)
	}
}

func (h *sseHub) broadcast(endpointID string, req *store.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	subs := h.subs[endpointID]
	if subs == nil {
		return
	}
	data, err := json.Marshal(req)
	if err != nil {
		return
	}
	for ch := range subs {
		select {
		case ch <- data:
		default:
		}
	}
}

// broadcastRequest publishes req to SSE subscribers for its endpoint.
func broadcastRequest(endpointID string, req *store.Request) {
	hub.broadcast(endpointID, req)
}

// HandleSSE serves an SSE stream for GET /e/{endpointID}/events.
func HandleSSE(repo store.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpointID := r.PathValue("endpointID")
		if endpointID == "" {
			http.Error(w, "missing endpoint ID", http.StatusBadRequest)
			return
		}
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)
		flusher.Flush()

		ch := hub.subscribe(endpointID)
		defer hub.unsubscribe(endpointID, ch)

		ctx := r.Context()
		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-ch:
				if !ok {
					return
				}
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
			}
		}
	}
}
