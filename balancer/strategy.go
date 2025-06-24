package balancer

import "github.com/abhi-shek-09/load-balancer/backend"

type Strategy interface {
	NextBackend([]*backend.Backend) *backend.Backend
}


// Any strategy must have a method NextBackend() 
// that takes the full list of backends and returns one of them