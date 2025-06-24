package balancer

import (
	"math/rand"
	"time"
	"github.com/abhi-shek-09/load-balancer/backend"
)

type Random struct {
	rng *rand.Rand
}

func NewRandom() *Random {
	seed := time.Now().UnixNano()
	return &Random{
		rng: rand.New(rand.NewSource(seed)),
	}
}

func (r *Random) NextBackend(backends []*backend.Backend) *backend.Backend {
	healthy := []*backend.Backend{}
	for _, b := range backends {
		if b.IsHealthy() {
			healthy = append(healthy, b)
		}
	}

	if len(healthy) == 0 {
		return nil
	}

	index := r.rng.Intn(len(healthy))
	return healthy[index]
}
