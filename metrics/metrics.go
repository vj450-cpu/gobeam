package metrics

import "gobeam/backend"

// BackendStat represents the stats for a single backend
type BackendStat struct {
	URL               string `json:"url"`
	ActiveConnections int    `json:"active"`
	RequestCount      int    `json:"requests"`
	Status            string `json:"status"`
}

// PoolStats represents the stats for the entire pool
type PoolStats struct {
	Servers []BackendStat `json:"servers"`
}

// GetPoolStats aggregates current stats from the pool
func GetPoolStats(pool *backend.ServerPool) PoolStats {
	var stats PoolStats
	for _, b := range pool.GetBackends() {
		active, requests, alive := b.GetStats()
		status := "DOWN"
		if alive {
			status = "UP"
		}
		stats.Servers = append(stats.Servers, BackendStat{
			URL:               b.URL.String(),
			Status:            status,
			ActiveConnections: active,
			RequestCount:      requests,
		})
	}
	return stats
}
