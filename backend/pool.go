package backend

import (
	"sync"
	"github.com/abhi-shek-09/load-balancer/config"
)

// Ensure the load balancer only routes to healthy backends, and automatically avoids routing to:
// Crashed servers
// Backends with no listener
// Temporarily unreachable instances

type Backend struct {
	Addr        string
	Healthy     bool
	ActiveConns int // for LeastConn
	Weight      int // for Weighted strategies
	mu          sync.RWMutex
}

var Backends []*Backend

func InitBackends() {
	for _, addr := range config.AppConfig.Backends {
		Backends = append(Backends, &Backend{Addr: addr, Healthy: true})
	}
}

func (b *Backend) SetHealth(healthy bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.Healthy = healthy
}

func (b *Backend) IsHealthy() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Healthy
}

func (b *Backend) IncConns() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.ActiveConns++
}

func (b *Backend) DecConns() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.ActiveConns > 0{
		b.ActiveConns--
	}
}

func (b *Backend) GetConns() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.ActiveConns
}

func FindBackendByAddr(addr string) *Backend {
	for _, b := range Backends {
		if b.Addr == addr {
			return b
		}
	}
	return nil
}
