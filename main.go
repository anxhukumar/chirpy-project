package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	currHits := cfg.fileserverHits.Load()
	formattedRes := fmt.Sprintf(`
<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>
`, currHits)
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

func handlerPostChirp(w http.ResponseWriter, r *http.Request) {
	// struct to receive a json
	type chirpData struct {
		Body string `json:"body"`
	}

	// struct to send a error json
	type errorData struct {
		Error string `json:"error"`
	}

	// struct to send a cleaned chirp json
	type cleanedMsg struct {
		CleanedBody string `json:"cleaned_body"`
	}

	// decode json
	decoder := json.NewDecoder(r.Body)
	chirp := chirpData{}
	err := decoder.Decode(&chirp)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Printf("Error decoding chirpData: %s", err)

		// send an error json in body
		errorResp := errorData{
			Error: "Something went wrong",
		}
		dat, err := json.Marshal(errorResp)
		if err != nil {
			log.Printf("Error marshalling error json: %s", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(500)
		w.Write(dat)
		return
	}

	// if the chirpData is longer than 140 characters
	if len(chirp.Body) > 140 {
		errorResp := errorData{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(errorResp)
		if err != nil {
			log.Printf("Error marshalling error json: %s", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	// send cleaned msg in body

	// update original chirp data
	chirp.Body = removeProfaneWords(chirp.Body)

	res := cleanedMsg{
		CleanedBody: chirp.Body,
	}

	dat, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling chirp json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(dat)
}

// helper functions
func removeProfaneWords(msg string) string {
	// split the words
	splitOriginalMsg := strings.Split(msg, " ")

	// lower case split
	splitLowercasedMsg := strings.Split(strings.ToLower(msg), " ")

	for i, v := range splitLowercasedMsg {
		if v == "kerfuffle" || v == "sharbert" || v == "fornax" {
			splitOriginalMsg[i] = "****"
		}
	}

	return strings.Join(splitOriginalMsg, " ")
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
