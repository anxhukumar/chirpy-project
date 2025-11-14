package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/anxhukumar/chirpy-project/internal/api"
	"github.com/anxhukumar/chirpy-project/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// middlewares
func MiddlewareMetricsInc(cfg *api.ApiConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
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
	apiConf := api.ApiConfig{
		FileserverHits: atomic.Int32{},
		Db:             dbQueries,
	}

	// ***add routing***
	mux.Handle("/app/", http.StripPrefix(
		"/app",
		MiddlewareMetricsInc(&apiConf, http.FileServer(http.Dir(filePathRoot))),
	))
	mux.HandleFunc("GET /api/healthz", api.HandlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiConf.HandlerRequestCount)
	mux.HandleFunc("POST /admin/reset", apiConf.HandlerReset)
	mux.HandleFunc("POST /api/users", apiConf.HandlerCreateUser)
	mux.HandleFunc("POST /api/chirps", apiConf.HandlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiConf.HandlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConf.HandlerGetOneChirp)

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
