package proxy

import (
	"io"
	"log"
	"net"
	"github.com/abhi-shek-09/load-balancer/balancer"
)

func HandleConnections(clientConn net.Conn) {
	defer clientConn.Close()

	backendAddr := balancer.GetNextBackend()
	if backendAddr == "" {
		log.Printf("No healthy backends available!")
		return
	}
	backendConn, err := net.Dial("tcp", backendAddr)
	if err != nil {
		log.Printf("Failed to connect to the backend %v", err)
		return
	}
	defer backendConn.Close()

	log.Printf("Proxying client %s <=> backend %s", clientConn.RemoteAddr(), clientConn.LocalAddr())

	// data piping
	go io.Copy(backendConn, clientConn) // client sending is a go routine so that other requests can be addressed
	io.Copy(clientConn, backendConn) // When this copy completes, we know the connection should terminate
	// TCP connections are bidirectional
	// both without goroutines, one would block and the other couldn't run
	// both with goroutines would get us out of the function, closing the connection, 
	// we'd need additional synchronization to keep the connection alive
}