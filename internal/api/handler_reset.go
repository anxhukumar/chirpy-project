package api

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func (cfg *ApiConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
	// check platform
	godotenv.Load()
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		log.Printf("Can't access this function without a local dev environment")
		return
	}
	// delete all users
	err := cfg.Db.DeleteAllUsers(r.Context())
	if err != nil {
		log.Printf("Error while resetting: %s", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	cfg.FileserverHits.Store(0)
}
