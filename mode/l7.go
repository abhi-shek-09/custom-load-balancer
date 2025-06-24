package mode

import (
	"github.com/abhi-shek-09/load-balancer/httppkg"
	"github.com/abhi-shek-09/load-balancer/security"
	"log"
	"net/http"
	"time"
)

func StartL7() {

	// re routing to secure traffic
	// runs a parallel server for redirection
	go func() {
		log.Println("🔁 Redirecting HTTP (8080) → HTTPS (8443)")
		err := http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			target := "https://" + r.Host + r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusPermanentRedirect) // 308
		}))
		if err != nil {
			log.Fatalf("Redirect server error: %v", err)
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/", security.IPFilterMiddleware(http.HandlerFunc(httppkg.ProxyHandler)))

	// SERVER TIMEOUT
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,  // Time to read request from client
		WriteTimeout: 10 * time.Second, // Time to write response to client
		IdleTimeout:  30 * time.Second, // Keep-alive timeout
	}

	log.Println("HTTP load balancer running on :8080 with timeouts")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// Timeout | 			Where				  |		What It Covers
// Server  | Load balancer's http.Server	  | Slow client reads, slow backend writes
// Backend | proxy transport (http.Transport) |  request delay, TCP connect time
