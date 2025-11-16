package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/anxhukumar/chirpy-project/internal/auth"
	"github.com/anxhukumar/chirpy-project/internal/database"
	"github.com/anxhukumar/chirpy-project/internal/helper"
)

func (cfg *ApiConfig) HandlerUpdateUser(w http.ResponseWriter, r *http.Request) {

	// fetch userId after validating JWT
	verifiedUserID := helper.FetchJwtUserId(cfg.JwtSecret, w, r)

	// fetch new changed user data
	newUserData := helper.DecodeAuthJson(w, r)

	// hash the password
	hashedPassword, err := auth.HashPassword(newUserData.Password)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// update data in database
	updateUserParams := database.UpdateUserByIdParams{
		Email:          newUserData.Email,
		HashedPassword: hashedPassword,
		ID:             verifiedUserID,
	}
	createdData, err := cfg.Db.UpdateUserById(r.Context(), updateUserParams)
	if err != nil {
		log.Printf("Error while updating user: %s", err)
		w.WriteHeader(500)
		return
	}

	// decode output
	userRes := struct {
		Id    string `json:"id"`
		Email string `json:"email"`
	}{
		Id:    createdData.Email,
		Email: createdData.Email,
	}
	res, err := json.Marshal(userRes)
	if err != nil {
		log.Printf("Error marshalling user data: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(res)
}
