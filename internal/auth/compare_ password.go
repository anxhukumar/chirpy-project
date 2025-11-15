package auth

import (
	"log"

	"github.com/alexedwards/argon2id"
)

// compares the inserted password with the actual password
func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Printf("Error checking password hash: %s", err)
		return false, err
	}
	return match, err
}
