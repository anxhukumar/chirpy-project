package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *ApiConfig) HandlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	// get all chirps
	chirps, err := cfg.Db.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Error while getting all chirps: %s", err)
		w.WriteHeader(500)
		return
	}

	// convert the slice of structs to slice of json responses

	// send chirp data one by one to json mapped struct instance
	newChirpSlice := []Chirp{}
	for _, v := range chirps {
		chirp := Chirp{
			ID:        v.ID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Body:      v.Body,
			UserID:    v.UserID,
		}

		newChirpSlice = append(newChirpSlice, chirp)
	}

	// convert the slice to json
	res, err := json.Marshal(newChirpSlice)
	if err != nil {
		log.Printf("Error converting slice of Chirps to json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(res)
}

func (cfg *ApiConfig) HandlerGetOneChirp(w http.ResponseWriter, r *http.Request) {
	// get id from url path
	idStr := r.PathValue("chirpID")

	// convert string id to uuid
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("Error converting ID string parameter to UUID: %s", err)
		w.WriteHeader(500)
		return
	}

	// get chirp
	chirp, err := cfg.Db.GetOneChirp(r.Context(), id)
	if err != nil {
		log.Printf("Error while getting one chirp: %s", err)
		w.WriteHeader(404)
		return
	}

	chirpData := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	// convert the struct to json
	res, err := json.Marshal(chirpData)
	if err != nil {
		log.Printf("Error converting chirp to json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(res)
}
