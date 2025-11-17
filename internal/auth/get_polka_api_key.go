package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	value := headers.Get("Authorization")
	if value == "" {
		log.Printf("Authorization header doesn't exist <get_polka_api_key>")
		return "", errors.New("auth header doesn't exist")
	}

	if !strings.HasPrefix(value, "ApiKey ") {
		log.Printf("ApiKey doesn't exist <get_polka_api_key>")
		return "", errors.New("ApiKey doesn't exist")
	}

	apiKey := strings.TrimPrefix(value, "ApiKey ")
	if apiKey == "" {
		log.Printf("ApiKey is empty  <get_polka_api_key>")
		return "", errors.New("ApiKey is empty")
	}

	return apiKey, nil
}
