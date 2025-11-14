package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// helper functions
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

func HandlerPostChirp(w http.ResponseWriter, r *http.Request) {
	// struct to receive a json
	type chirpData struct {
		Body string `json:"body"`
	}

	// struct to send a error json
	type errorData struct {
		Error string `json:"error"`
	}

	// struct to send a cleaned chirp json
	type cleanedMsg struct {
		CleanedBody string `json:"cleaned_body"`
	}

	// decode json
	decoder := json.NewDecoder(r.Body)
	chirp := chirpData{}
	err := decoder.Decode(&chirp)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Printf("Error decoding chirpData: %s", err)

		// send an error json in body
		errorResp := errorData{
			Error: "Something went wrong",
		}
		dat, err := json.Marshal(errorResp)
		if err != nil {
			log.Printf("Error marshalling error json: %s", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(500)
		w.Write(dat)
		return
	}

	// if the chirpData is longer than 140 characters
	if len(chirp.Body) > 140 {
		errorResp := errorData{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(errorResp)
		if err != nil {
			log.Printf("Error marshalling error json: %s", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	// send cleaned msg in body

	// update original chirp data
	chirp.Body = RemoveProfaneWords(chirp.Body)

	res := cleanedMsg{
		CleanedBody: chirp.Body,
	}

	dat, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling chirp json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(dat)
}
