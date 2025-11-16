package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
)

func MakeRefreshToken() (string, error) {
	// generate random code
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Printf("Error while creating random hex string in MakeRefreshToken(): %s", err)
		return "", errors.New("error executing rand.Read(key) in MakeRefreshToken()")
	}

	encoded_str := hex.EncodeToString(key)
	return encoded_str, nil
}
