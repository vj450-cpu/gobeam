package dashboard

import (
	"gobeam/backend"
	"log"
	"net/http"
)

// StartServer starts the dashboard HTTP server
func StartServer(port string, pool *backend.ServerPool) {
	mux := http.NewServeMux()
	
	// Serve the static HTML dashboard
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "dashboard/index.html")
	})
	
	// API endpoint
	mux.HandleFunc("/api/stats", APIHandler(pool))
	
	log.Printf("Starting dashboard on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Printf("Dashboard server stopped: %v", err)
	}
}
