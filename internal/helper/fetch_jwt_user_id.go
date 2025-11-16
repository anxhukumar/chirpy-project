package helper

import (
	"net/http"

	"github.com/anxhukumar/chirpy-project/internal/auth"
	"github.com/google/uuid"
)

func FetchJwtUserId(jwtSecret string, w http.ResponseWriter, r *http.Request) uuid.UUID {
	// fetch access token
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return uuid.Nil
	}

	// validate received token
	verifiedUserID, err := auth.ValidateJWT(accessToken, jwtSecret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return uuid.Nil
	}

	return verifiedUserID
}
