package loadbalancer

import (
	"gobeam/backend"
	"sync"
)

// RoundRobin implements weighted round robin load balancing
type RoundRobin struct {
	mux           sync.Mutex
	currentIdx    int
	currentWeight int
}

// NextServer selects the next available server using weighted round robin
func (r *RoundRobin) NextServer(pool *backend.ServerPool) *backend.Backend {
	backends := pool.GetBackends()
	n := len(backends)
	if n == 0 {
		return nil
	}

	r.mux.Lock()
	defer r.mux.Unlock()

	// Find the max weight in the pool
	maxWeight := 0
	for _, b := range backends {
		if b.Weight > maxWeight {
			maxWeight = b.Weight
		}
	}
	if maxWeight == 0 {
		return nil
	}

	for {
		r.currentIdx = (r.currentIdx + 1) % n
		if r.currentIdx == 0 {
			r.currentWeight--
			if r.currentWeight <= 0 {
				r.currentWeight = maxWeight
				if r.currentWeight == 0 {
					return nil
				}
			}
		}

		b := backends[r.currentIdx]
		if b.IsAlive() && b.Weight >= r.currentWeight {
			return b
		}
	}
}
