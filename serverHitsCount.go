package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/anxhukumar/chirpy-project/internal/database"
)

// serverHits config
type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

// resets the server hits
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
}

// returns current hits
func (cfg *apiConfig) handlerRequestCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	currHits := cfg.fileserverHits.Load()
	formattedRes := fmt.Sprintf(`
<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>
`, currHits)
	w.Write([]byte(formattedRes))
}
