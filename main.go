// Layer 4 load balancing works at the transport layer, focusing on IP addresses and ports, 
// while Layer 7 load balancing operates at the application layer, considering the content of the application data (like HTTP headers) to make routing decisions. 
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"github.com/abhi-shek-09/load-balancer/health"
	"github.com/abhi-shek-09/load-balancer/mode"
	"github.com/abhi-shek-09/load-balancer/config"
	"github.com/abhi-shek-09/load-balancer/balancer"
	"github.com/abhi-shek-09/load-balancer/backend"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to config.yaml")
	flag.Parse()
	
	config.LoadConfig(*configPath)
	backend.InitBackends()
	balancer.InitStrategy()
	health.StartHealthCheck()

	switch config.AppConfig.Mode {
	case "l4":
		log.Println("Running in L4 TCP mode")
		mode.StartL4()

	case "l7":
		log.Println("Running in L7 HTTP mode")
		switch strings.ToLower(config.AppConfig.Strategy) {
		case "roundrobin":
			balancer.SetStrategy(balancer.NewRoundRobin())
		case "random":
			balancer.SetStrategy(balancer.NewRandom())
		case "leastconn":
			balancer.SetStrategy(balancer.NewLeastConnections())
		default:
			log.Fatalf("Unknown strategy: %s", config.AppConfig.Strategy)
		}

		fmt.Printf("🔁 Strategy set to: %s\n", config.AppConfig.Strategy)

		if config.AppConfig.TLS {
			mode.StartTLSL7()
		} else {
			mode.StartL7()
		}
		mode.StartL7()

	default:
		log.Fatal("Unknown or missing --mode. Use 'l4' or 'l7'")
	}
}
