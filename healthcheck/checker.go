package healthcheck

import (
	"gobeam/backend"
	"log"
	"net/http"
	"time"
)

// StartHealthCheck starts a background goroutine to check backend health
func StartHealthCheck(pool *backend.ServerPool, interval time.Duration) {
	go func() {
		t := time.NewTicker(interval)
		defer t.Stop()

		for {
			<-t.C
			for _, b := range pool.GetBackends() {
				// Time out quickly since we are checking health frequently
				client := http.Client{
					Timeout: 2 * time.Second,
				}
				healthURL := b.URL.String() + "/health"
				resp, err := client.Get(healthURL)

				status := "UP"
				isAlive := true

				if err != nil || resp.StatusCode != http.StatusOK {
					isAlive = false
					status = "DOWN"
				}

				if resp != nil {
					resp.Body.Close()
				}

				wasAlive := b.IsAlive()
				if wasAlive != isAlive {
					log.Printf("Backend state change for %s. Status: %s", b.URL, status)
				}
				
				b.SetAlive(isAlive)
			}
		}
	}()
}
