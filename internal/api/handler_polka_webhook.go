package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/anxhukumar/chirpy-project/internal/auth"
	"github.com/google/uuid"
)

type PolkaReq struct {
	Event string `json:"event"`
	Data  struct {
		UserId string `json:"user_id"`
	} `json:"data"`
}

func (cfg *ApiConfig) HandlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {

	//verify if request is coming from correct source
	receivedApiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if receivedApiKey != cfg.PolkaApiKey {
		log.Printf("Polka webhook request made from unauthorized source")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// decode polka json data
	decoder := json.NewDecoder(r.Body)
	polkaData := PolkaReq{}
	err = decoder.Decode(&polkaData)
	if err != nil {
		log.Printf("Error decoding polka request: %s", err)
		w.WriteHeader(500)
		return
	}

	// check the event type
	if polkaData.Event != "user.upgraded" {
		log.Printf("Polka sent an irrelevant request event: %s", polkaData.Event)
		w.WriteHeader(204)
		return
	}

	// convert string id to uuid
	userId, err := uuid.Parse(polkaData.Data.UserId)
	if err != nil {
		log.Printf("Error converting userId string parameter to UUID in <handler_polka_webhook>: %s", err)
		w.WriteHeader(500)
		return
	}

	// update is_chirpy_red in database
	err = cfg.Db.UpdateIsChirpyRedById(r.Context(), userId)
	if err != nil {
		log.Printf("Error updating is_chirpy_red or user not found in <handler_polka_webhook>: %s", err)
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(204)
}
