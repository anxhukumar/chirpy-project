package helper

import "strings"

func RemoveProfaneWords(msg string) string {
	// split the words
	splitOriginalMsg := strings.Split(msg, " ")

	// lower case split
	splitLowercasedMsg := strings.Split(strings.ToLower(msg), " ")

	// bad words map
	badWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	// iterate and replace the bad words if they exist
	for i, v := range splitLowercasedMsg {
		if badWords[v] {
			splitOriginalMsg[i] = "****"
		}
	}

	return strings.Join(splitOriginalMsg, " ")
}
