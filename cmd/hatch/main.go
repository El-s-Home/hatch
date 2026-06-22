package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/elfoundation/hatch/internal/handler"
	"github.com/elfoundation/hatch/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	dbPath := os.Getenv("HATCH_DB_PATH")
	repo, err := store.Open(dbPath)
	if err != nil { log.Fatalf("hatch: open store: %v", err) }
	defer repo.Close()
	h := handler.New(repo)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthz)
	h.RegisterRoutes(mux)
	addr := fmt.Sprintf(":%s", port)
	log.Printf("hatch starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil { log.Fatalf("hatch server error: %v", err) }
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}
