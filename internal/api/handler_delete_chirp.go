package api

import (
	"log"
	"net/http"

	"github.com/anxhukumar/chirpy-project/internal/helper"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandlerDeleteChirp(w http.ResponseWriter, r *http.Request) {

	// fetch userId after validating JWT
	verifiedUserID := helper.FetchJwtUserId(cfg.JwtSecret, w, r)

	// get id from url path
	chirpIdStr := r.PathValue("chirpID")

	// convert string id to uuid
	chirpId, err := uuid.Parse(chirpIdStr)
	if err != nil {
		log.Printf("Error converting ID string parameter to UUID: %s", err)
		w.WriteHeader(500)
		return
	}

	// get chirp
	chirp, err := cfg.Db.GetOneChirp(r.Context(), chirpId)
	if err != nil {
		log.Printf("Error while getting chirp data | <handler_delete_chirp>: %s", err)
		w.WriteHeader(404)
		return
	}

	// make sure the the userId of chirp and verified user id is same
	if chirp.UserID != verifiedUserID {
		log.Printf("User id of access token != user_id of chirp | <handler_delete_chirp>: %s", err)
		w.WriteHeader(403)
		return
	}

	// delete the chirp
	err = cfg.Db.DeleteOneChirp(r.Context(), chirp.ID)
	if err != nil {
		log.Printf("Error while deleting chirp: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)

}
