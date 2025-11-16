package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/anxhukumar/chirpy-project/internal/auth"
	"github.com/anxhukumar/chirpy-project/internal/database"
	"github.com/anxhukumar/chirpy-project/internal/helper"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// decode user auth json
	userData := helper.DecodeAuthJson(w, r)

	// create user in database

	// hash the password
	hashedPassword, err := auth.HashPassword(userData.Password)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	userAuthData := database.CreateUserParams{
		Email:          userData.Email,
		HashedPassword: hashedPassword,
	}
	createdData, err := cfg.Db.CreateUser(r.Context(), userAuthData)
	if err != nil {
		log.Printf("Couldn't create user: %s", err)
		w.WriteHeader(500)
		return
	}

	// decode output
	userRes := User{
		ID:        createdData.ID,
		CreatedAt: createdData.CreatedAt,
		UpdatedAt: createdData.UpdatedAt,
		Email:     createdData.Email,
	}
	res, err := json.Marshal(userRes)
	if err != nil {
		log.Printf("Error marshalling user data: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(res)

}
