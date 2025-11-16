package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/anxhukumar/chirpy-project/internal/auth"
	"github.com/anxhukumar/chirpy-project/internal/database"
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
	jwtToken, err := helper.GetJwtToken(userData, cfg.JwtSecret)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// get refresh token
	refreshTokenCode, err := auth.MakeRefreshToken()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	// send the refresh token to database
	REFRESH_TOKEN_EXPIRY_DURATION_IN_DAYS := 60
	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refreshTokenCode,
		UserID:    userData.ID,
		ExpiresAt: time.Now().Add(time.Duration(REFRESH_TOKEN_EXPIRY_DURATION_IN_DAYS*24) * time.Hour),
	}
	refreshToken, err := cfg.Db.CreateRefreshToken(r.Context(), refreshTokenParams)
	if err != nil {
		log.Printf("Error while inserting refreshToken in database: %s", err)
		return
	}

	// return user data since the password has matched
	userRes := User{
		ID:           userData.ID,
		CreatedAt:    userData.CreatedAt,
		UpdatedAt:    userData.UpdatedAt,
		Email:        userData.Email.String,
		Token:        jwtToken,
		RefreshToken: refreshToken.Token,
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
