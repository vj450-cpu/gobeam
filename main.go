package main

import (
	"encoding/json"
	"gobeam/backend"
	"gobeam/dashboard"
	"gobeam/healthcheck"
	"gobeam/loadbalancer"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync/atomic"
	"time"
)

type ServerConfig struct {
	URL    string `json:"url"`
	Weight int    `json:"weight"`
}

type Config struct {
	Servers []ServerConfig `json:"servers"`
}

func main() {
	// Read config
	configFile, err := os.ReadFile("config/servers.json")
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	pool := &backend.ServerPool{}

	// Initialize backends
	for _, s := range config.Servers {
		parsedUrl, err := url.Parse(s.URL)
		if err != nil {
			log.Fatalf("Invalid server URL %s: %v", s.URL, err)
		}

		proxy := httputil.NewSingleHostReverseProxy(parsedUrl)

		b := &backend.Backend{
			URL:          parsedUrl,
			Alive:        true,
			Weight:       s.Weight,
			ReverseProxy: proxy,
		}

		// Setup ErrorHandler to mark server as dead on proxy error
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
			log.Printf("[%s] Proxy error: %v", b.URL.Host, e)
			b.SetAlive(false)
			http.Error(w, "Backend unavailable", http.StatusServiceUnavailable)
		}

		pool.AddBackend(b)
		log.Printf("Configured backend: %s (weight: %d)", parsedUrl, s.Weight)
	}

	// Setting up balancer algorithm (Weighted Round Robin)
	balancer := &loadbalancer.RoundRobin{}

	// Global rate limiter setup
	var requestCounter int32
	go func() {
		for {
			time.Sleep(1 * time.Second)
			atomic.StoreInt32(&requestCounter, 0)
		}
	}()

	// Create main server proxy logic
	server := http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if atomic.AddInt32(&requestCounter, 1) > 50 {
				log.Printf("[RateLimiter] Request blocked from %s", r.RemoteAddr)
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			b := balancer.NextServer(pool)
			if b == nil {
				http.Error(w, "Service not available", http.StatusServiceUnavailable)
				return
			}

			b.IncrementConnections()
			defer b.DecrementConnections()

			// Proxy Request
			b.ReverseProxy.ServeHTTP(w, r)
		}),
	}

	// Start health checking background process
	healthcheck.StartHealthCheck(pool, 5*time.Second)

	// Start dashboard on port 9000 in background
	go dashboard.StartServer("9000", pool)

	// Start load balancer
	log.Println("Starting Load Balancer on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
