package loadbalancer

import "gobeam/backend"

// Balancer is the interface for load balancing algorithms
type Balancer interface {
	NextServer(pool *backend.ServerPool) *backend.Backend
}
