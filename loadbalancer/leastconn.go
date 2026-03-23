package loadbalancer

import "gobeam/backend"

// LeastConn implements least connections load balancing
type LeastConn struct{}

// NextServer selects the alive server with the fewest active connections
func (l *LeastConn) NextServer(pool *backend.ServerPool) *backend.Backend {
	backends := pool.GetBackends()
	var best *backend.Backend
	var bestConn int
	
	for _, b := range backends {
		if !b.IsAlive() {
			continue
		}
		
		bConn, _, _ := b.GetStats()
		
		if best == nil {
			best = b
			bestConn = bConn
			continue
		}
		
		if bConn < bestConn {
			best = b
			bestConn = bConn
		}
	}
	
	return best
}
