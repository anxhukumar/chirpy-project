package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/anxhukumar/chirpy-project/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
	// load env
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	// connect sql db
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	// port config
	const filePathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	// apiconfig
	apiConfig := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
	}

	// add routing
	mux.Handle("/app/", http.StripPrefix(
		"/app",
		apiConfig.middlewareMetricsInc(http.FileServer(http.Dir(filePathRoot))),
	))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiConfig.handlerRequestCount)
	mux.HandleFunc("POST /admin/reset", apiConfig.handlerReset)
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
