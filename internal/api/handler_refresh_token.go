package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/anxhukumar/chirpy-project/internal/auth"
	"github.com/anxhukumar/chirpy-project/internal/helper"
)

type NewToken struct {
	Token string `json:"token"`
}

func (cfg *ApiConfig) HandlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	// get refresh token from client
	refreshTokenID, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// create a fresh access token
	userData, err := cfg.Db.GetUserFromRefreshToken(r.Context(), refreshTokenID)
	if err != nil {
		log.Printf("Error while getting user from refresh token in HandlerRefreshToken or it is revoked: %s", err)
		w.WriteHeader(401)
		return
	}

	// get jwt token
	jwtToken, err := helper.GetJwtToken(userData, cfg.JwtSecret)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// send token json res
	tokenRes := NewToken{
		Token: jwtToken,
	}
	res, err := json.Marshal(tokenRes)
	if err != nil {
		log.Printf("Error marshalling new token data: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(res)
}
