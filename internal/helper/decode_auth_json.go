package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

type UserAuthData struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

func DecodeAuthJson(w http.ResponseWriter, r *http.Request) UserAuthData {
	decoder := json.NewDecoder(r.Body)
	userData := UserAuthData{}
	err := decoder.Decode(&userData)
	if err != nil {
		log.Printf("Error decoding emailData: %s", err)
		w.WriteHeader(500)
		return UserAuthData{}
	}
	return userData
}
