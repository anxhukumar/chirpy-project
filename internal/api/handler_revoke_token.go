package api

import (
	"log"
	"net/http"

	"github.com/anxhukumar/chirpy-project/internal/auth"
)

func (cfg *ApiConfig) HandlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	// get refresh token from client
	refreshTokenID, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// set revoked_at time value in refresh_token table in database
	err = cfg.Db.UpdateRevokedAt(r.Context(), refreshTokenID)
	if err != nil {
		log.Printf("Error while updating revoked_at timestamp: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}
