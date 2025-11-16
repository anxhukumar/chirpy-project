package helper

import (
	"time"

	"github.com/anxhukumar/chirpy-project/internal/auth"
	"github.com/anxhukumar/chirpy-project/internal/database"
)

const JWT_TOKEN_EXPIRY_DURATION_IN_SECONDS = 3600

func GetJwtToken(userData database.User, jwtSecret string) (string, error) {
	jwtToken, err := auth.MakeJWT(
		userData.ID,
		jwtSecret,
		time.Duration(JWT_TOKEN_EXPIRY_DURATION_IN_SECONDS)*time.Second,
	)
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}
