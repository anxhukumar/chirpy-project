package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/anxhukumar/chirpy-project/internal/auth"
	"github.com/anxhukumar/chirpy-project/internal/helper"
)

func (cfg *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	// decode user auth json
	loginData := helper.DecodeAuthJson(w, r)

	// check and get if user exists
	userData, err := cfg.Db.GetUserByEmail(r.Context(), helper.ToNullString(loginData.Email))
	if err != nil {
		log.Printf("Incorrect email or password: %s", err)
		w.WriteHeader(401)
		return
	}

	// compare password
	isMatch, err := auth.CheckPasswordHash(loginData.Password, userData.HashedPassword)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if !isMatch {
		log.Printf("Incorrect email or password: %s", err)
		w.WriteHeader(401)
		return
	}

	// get jwt token
	if loginData.ExpiresInSeconds == 0 || loginData.ExpiresInSeconds > 3600 {
		loginData.ExpiresInSeconds = 3600
	}
	jwtToken, err := auth.MakeJWT(userData.ID, cfg.JwtSecret, time.Duration(loginData.ExpiresInSeconds)*time.Second)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// return user data since the password has matched
	userRes := User{
		ID:        userData.ID,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
		Email:     userData.Email.String,
		Token:     jwtToken,
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
