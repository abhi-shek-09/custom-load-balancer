package mode

import (
	"log"
	"net"
	"github.com/abhi-shek-09/load-balancer/proxy"
)

func StartL4() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error starting TCP listener: %v", err)
	}
	defer listener.Close()

	log.Println("Load balancer listening on :8080")

	for {
		// Accept() is initially waiting for a connection
		// Accept() wakes up and returns a new net.Conn
		// After that it goes back to sleep, waiting for a connection, coz its a go routine
		// if it was proxy.HandleConnections(conn), we could have had only 1 connection
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept: %v", err)
			continue
		}
		go proxy.HandleConnections(conn)
	}
}