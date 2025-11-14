package api

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/anxhukumar/chirpy-project/internal/database"
)

// serverHits config
type ApiConfig struct {
	FileserverHits atomic.Int32
	Db             *database.Queries
}

// returns current hits
func (cfg *ApiConfig) HandlerRequestCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	currHits := cfg.FileserverHits.Load()
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
