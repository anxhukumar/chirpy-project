package helper

import (
	"log"
	"net/http"
)

func ValidateChirp(chirp string, w http.ResponseWriter) string {

	if len(chirp) > 140 {
		log.Printf("The length of chirp is too long")
		w.WriteHeader(413)
		return ""
	}

	// filter profanity from chirp
	cleanChirp := RemoveProfaneWords(chirp)

	return cleanChirp
}
