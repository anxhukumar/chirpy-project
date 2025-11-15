package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	value := headers.Get("Authorization")
	if value == "" {
		log.Printf("Authorization header doesn't exist")
		return "", errors.New("auth header doesn't exist")
	}

	if !strings.HasPrefix(value, "Bearer ") {
		log.Printf("Bearer token doesn't exist")
		return "", errors.New("bearer token doesn't exist")
	}

	token := strings.TrimPrefix(value, "Bearer ")
	if token == "" {
		log.Printf("Bearer token is empty")
		return "", errors.New("bearer token is empty")
	}

	return token, nil
}
