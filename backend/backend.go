package backend

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

// Backend holds the data about a server
type Backend struct {
	URL               *url.URL
	Alive             bool
	mux               sync.RWMutex
	ActiveConnections int
	RequestCount      int
	Weight            int
	ReverseProxy      *httputil.ReverseProxy
}

// SetAlive updates the Alive status
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

// IsAlive returns true when backend is alive
func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	alive := b.Alive
	b.mux.RUnlock()
	return alive
}

// IncrementConnections increments active connections and request count
func (b *Backend) IncrementConnections() {
	b.mux.Lock()
	b.ActiveConnections++
	b.RequestCount++
	b.mux.Unlock()
}

// DecrementConnections decrements active connections
func (b *Backend) DecrementConnections() {
	b.mux.Lock()
	b.ActiveConnections--
	b.mux.Unlock()
}

// GetStats returns current stats for dashboard
func (b *Backend) GetStats() (int, int, bool) {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.ActiveConnections, b.RequestCount, b.Alive
}
