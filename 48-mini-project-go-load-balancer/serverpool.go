package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

// ServerPool holds information about reachable backends
type ServerPool struct {
	backends []*Backend
	current  uint64
}

// AddBackend to the server pool
func (s *ServerPool) AddBackend(backend *Backend) {
	s.backends = append(s.backends, backend)
}

// NextIndex atomically increases the counter and returns an index
func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

// GetNextPeer returns next active backend to take a connection
func (s *ServerPool) GetNextPeer() *Backend {
	// loop entire backends to find an alive one
	next := s.NextIndex()
	l := len(s.backends) + next // start from next and cycle a whole loop if needed
	for i := next; i < l; i++ {
		idx := i % len(s.backends)
		// if backend is alive and isn't the current one, return it
		if s.backends[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx))
			}
			return s.backends[idx]
		}
	}
	return nil
}

// HealthCheck pings the backends and update the status
func (s *ServerPool) HealthCheck() {
	for _, b := range s.backends {
		alive := isBackendAlive(b.URL)
		b.SetAlive(alive)
		if !alive {
			log.Printf("%s is down\n", b.URL)
		}
	}
}

// LB load balances the incoming request
func (s *ServerPool) LB(w http.ResponseWriter, r *http.Request) {
	peer := s.GetNextPeer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}
