package dashboard

import (
	"encoding/json"
	"gobeam/backend"
	"gobeam/metrics"
	"net/http"
)

// APIHandler handles the /api/stats endpoint
func APIHandler(pool *backend.ServerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := metrics.GetPoolStats(pool)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}
