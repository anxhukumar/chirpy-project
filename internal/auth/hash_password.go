package auth

import (
	"log"

	"github.com/alexedwards/argon2id"
)

// hash the password using argon2id
func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		return "", err
	}

	return hash, err
}
