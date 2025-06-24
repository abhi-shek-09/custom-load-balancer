package balancer

import (
	"log"
	"math"

	"github.com/abhi-shek-09/load-balancer/backend"
)

type LeastConnections struct{}

func NewLeastConnections() *LeastConnections{
	return &LeastConnections{}
}

func (lc *LeastConnections) NextBackend(backends []*backend.Backend) *backend.Backend {
	var selected *backend.Backend
	minConns := math.MaxInt32

	for _, b := range backends {
		if !b.IsHealthy(){
			continue
		}

		conns := b.ActiveConns
		log.Println(conns)
		if conns < minConns {
			minConns = conns
			selected = b
		}
	}

	selected.IncConns()
	return selected
}