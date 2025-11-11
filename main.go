package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

// serverHits config
type apiConfig struct {
	fileserverHits atomic.Int32
}

// methods of serverHits config
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// returns current hits
func (cfg *apiConfig) handlerRequestCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	currHits := cfg.fileserverHits.Load()
	formattedRes := fmt.Sprintf("Hits: %v\n", currHits)
	w.Write([]byte(formattedRes))
}

// resets the server hits
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
}

// Non-method handlers
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {
	// port config
	const filePathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	// serverHitsCount state
	serveHitsCount := apiConfig{}

	// add routing
	mux.Handle("/app/", http.StripPrefix(
		"/app",
		serveHitsCount.middlewareMetricsInc(http.FileServer(http.Dir(filePathRoot))),
	))
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", serveHitsCount.handlerRequestCount)
	mux.HandleFunc("/reset", serveHitsCount.handlerReset)

	// server struct
	serv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	// log confirmation that the port is running
	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)

	// if the server fails it will log the error
	log.Fatal(serv.ListenAndServe())
}
