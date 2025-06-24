package balancer

import (
	"sync"
	"github.com/abhi-shek-09/load-balancer/backend"
)

type RoundRobin struct {
	current int
	mu sync.Mutex
}

// Encapsulates creation logic
// All it is doing is creating a round robin strategy and giving it
// We want to keep track of which backend was picked last time
func NewRoundRobin() *RoundRobin {
	return &RoundRobin{}
}

// We will be writing NextBackend for every strategy
// I want this method to work on my specific RoundRobin object (rr)
// coz we want to avoid copying the whole backend struct since thats inefficient
// also some fields inside like Healthy use mutexes and shld stay safe and shared
// if u dont use pointers and by any chance copy the mutex, your program will break

func (rr *RoundRobin) NextBackend(backends []*backend.Backend) *backend.Backend {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	start := rr.current
	for {
		bEnd := backends[rr.current]
		rr.current = (rr.current+1) % len(backends)

		if bEnd.IsHealthy() {
			return bEnd
		}

		if start == rr.current {
			break
		}
	}
	return nil
}