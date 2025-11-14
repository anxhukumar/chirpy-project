package main

import (
	"log"
	"net/http"
)

// middlewares

// method of serverHits config
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
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
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", serveHitsCount.handlerRequestCount)
	mux.HandleFunc("POST /admin/reset", serveHitsCount.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerPostChirp)

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
