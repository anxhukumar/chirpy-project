package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/anxhukumar/chirpy-project/internal/helper"
	"github.com/google/uuid"
)

type Email struct {
	Email string `json:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// decode email json
	decoder := json.NewDecoder(r.Body)
	email := Email{}
	err := decoder.Decode(&email)
	if err != nil {
		log.Printf("Error decoding emailData: %s", err)
		w.WriteHeader(500)
		return
	}

	// create user in database
	createdData, err := cfg.Db.CreateUser(r.Context(), helper.ToNullString(email.Email))
	if err != nil {
		log.Printf("Couldn't create user: %s", err)
		w.WriteHeader(500)
		return
	}

	// decode output
	userData := User{
		ID:        createdData.ID,
		CreatedAt: createdData.CreatedAt,
		UpdatedAt: createdData.UpdatedAt,
		Email:     createdData.Email.String,
	}
	res, err := json.Marshal(userData)
	if err != nil {
		log.Printf("Error marshalling user data: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(res)

}
