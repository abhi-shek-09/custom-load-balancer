package mode

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"
	"github.com/abhi-shek-09/load-balancer/httppkg"
	"github.com/abhi-shek-09/load-balancer/security"
)

func StartTLSL7() {

	go func() {
		log.Println("Redirecting 8080 to 8443")
		err := http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			target := "https://" + r.Host + r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusPermanentRedirect)
		}))
		if err != nil {
			log.Fatalf("Redirect server error : %v", err)
		}
	}()
	
	mux := http.NewServeMux()
	mux.Handle("/", security.RateLimitMiddleware(security.IPFilterMiddleware(http.HandlerFunc(httppkg.ProxyHandler))))

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		// You can tweak cipher suites, ALPN, etc. later
	}

	server := &http.Server{
		Addr:         ":8443",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
		TLSConfig:    tlsConfig,
	}

	log.Println("🔒 HTTPS load balancer running on :8443")
	err := server.ListenAndServeTLS("cert/server.crt", "cert/server.key")
	if err != nil {
		log.Fatal("TLS Server error:", err)
	}

}
