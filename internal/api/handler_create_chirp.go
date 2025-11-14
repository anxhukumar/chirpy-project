package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/anxhukumar/chirpy-project/internal/database"
	"github.com/anxhukumar/chirpy-project/internal/helper"
	"github.com/google/uuid"
)

type ChirpData struct {
	Body   string        `json:"body"`
	UserID uuid.NullUUID `json:"user_id"`
}

type ChirpFullData struct {
	ID        uuid.UUID     `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Body      string        `json:"body"`
	UserID    uuid.NullUUID `json:"user_id"`
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

	// get validated chirp
	// return invalid code to user from ValidateChirp function if the body is invalid
	validChirp := helper.ValidateChirp(chirpData.Body, w)

	// create user in database
	createdData, err := cfg.Db.CreateChirp(r.Context(),
		database.CreateChirpParams{Body: validChirp, UserID: chirpData.UserID})
	if err != nil {
		log.Printf("Couldn't create chirp: %s", err)
		w.WriteHeader(500)
		return
	}

	//  decode output
	ChirpFullData := ChirpFullData{
		ID:        createdData.ID,
		CreatedAt: createdData.CreatedAt,
		UpdatedAt: createdData.UpdatedAt,
		Body:      createdData.Body,
		UserID:    createdData.UserID,
	}
	res, err := json.Marshal(ChirpFullData)
	if err != nil {
		log.Printf("Error marshalling chirp data: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(res)

}
