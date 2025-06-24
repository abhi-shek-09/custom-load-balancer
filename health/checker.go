package health

import (
	"log"
	"net"
	"time"
	"github.com/abhi-shek-09/load-balancer/backend"
)

const interval = 5 * time.Second
const timeout = 2 * time.Second

func StartHealthCheck() {
	go func() {
		for {
			for _, b := range backend.Backends {
				go checkBackend(b)
			}
			time.Sleep(interval)
		}
	} ()
}

func checkBackend(b *backend.Backend) {
	// log only if theres a change in behavior
	conn, err := net.DialTimeout("tcp", b.Addr, timeout)
	wasHealthy := b.IsHealthy()

	if err != nil {
		b.SetHealth(false)
		// it was healthy, but now conn failed, so make it unhealthy and report
		if wasHealthy {
			log.Printf("%s is marked Unhealthy", b.Addr)
		}
		return
	}

	conn.Close()
	b.SetHealth(true)

	// was unhealthy, but now conn passed, so make it healthy and report
	if !wasHealthy {
		log.Printf("%s marked Healthy", b.Addr)
	}
}

// Starts automatically in background when StartHealthCheck() is called

// Every 5 seconds, it:
	// Loops through all registered backends
	// Spawns a goroutine for each backend check (parallel checking)

// Each check:
	// Attempts TCP connection with 2-second timeout
	// If connection fails → marks backend unhealthy
	// If connection succeeds → marks backend healthy
	// Only logs status changes (not every check)