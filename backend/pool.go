package backend

// ServerPool holds a list of backends
type ServerPool struct {
	Backends []*Backend
}

// AddBackend adds a new backend to the pool
func (s *ServerPool) AddBackend(b *Backend) {
	s.Backends = append(s.Backends, b)
}

// GetBackends returns all backends in the pool
func (s *ServerPool) GetBackends() []*Backend {
	return s.Backends
}
