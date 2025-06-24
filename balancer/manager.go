package balancer

import (
	"github.com/abhi-shek-09/load-balancer/backend"
	"github.com/abhi-shek-09/load-balancer/config"
	"strings"
	"sync"
)

var (
	mu              sync.RWMutex
	currentStrategy Strategy
)

func InitStrategy() {
	switch strings.ToLower(config.AppConfig.Strategy) {
	case "least_connections":
		currentStrategy = NewLeastConnections()
	case "random":
		currentStrategy = NewRandom()
	default:
		currentStrategy = NewRoundRobin()
	}
}

func SetStrategy(s Strategy) {
	mu.Lock()
	defer mu.Unlock()
	currentStrategy = s
}

func GetNextBackend() string {
	mu.RLock()
	defer mu.RUnlock()

	if b := currentStrategy.NextBackend(backend.Backends); b != nil {
		return b.Addr
	}
	return ""
}
