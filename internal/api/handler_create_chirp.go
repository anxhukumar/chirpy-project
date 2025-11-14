package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type ChirpData struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}

func (cfg *ApiConfig) HandlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	// decode chirp input
	decoder := json.NewDecoder(r.Body)
	chirpData := ChirpData{}
	err := decoder.Decode(&chirpData)
	if err != nil {
		log.Printf("Error decoding chirp data: %s", err)
		w.WriteHeader(500)
		return
	}
}
