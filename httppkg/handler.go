package httppkg

import (
	"log"
	"net"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"time"
	"github.com/abhi-shek-09/load-balancer/balancer"
	"github.com/abhi-shek-09/load-balancer/backend"
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	
	const maxRetries = 1
	attempts := 0
	start :=  time.Now()

	if r.TLS != nil {
		log.Printf("🔐 TLS Info: Version=%x, CipherSuite=%x", r.TLS.Version, r.TLS.CipherSuite)
	}
	
	for {
		backendAddr := balancer.GetNextBackend()
		if backendAddr == "" {
			http.Error(w, "No healthy backend available", http.StatusServiceUnavailable)
			return
		}

		target, err := url.Parse("http://" + backendAddr)
		if err != nil {
			http.Error(w, "Invalid backend address", http.StatusInternalServerError)
			return
		}
		
		log.Printf("HTTP %s %s → %s", r.Method, r.URL.Path, backendAddr)
		// A proxy server (or forward proxy) acts on behalf of a client, hiding the clients identity and potentially filtering requests.
		// A reverse proxy, on the other hand, acts on behalf of a server, accepting requests from clients and forwarding them to the appropriate backend servers
		proxy := httputil.NewSingleHostReverseProxy(target)

		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Host = r.Host
		}

		// Backend timeout
		proxy.Transport = &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   3 * time.Second, // tcp connection timeout to backend
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ResponseHeaderTimeout: 5 * time.Second,
			IdleConnTimeout:       90 * time.Second,
			MaxIdleConns:          100,
		}

		// record the response
		rec := httptest.NewRecorder()
		b := backend.FindBackendByAddr(backendAddr)
		if b != nil {
			b.IncConns()
			defer b.DecConns() // ensure it decreases no matter what
		}
		proxy.ServeHTTP(rec, r)

		status := rec.Result().StatusCode
		if status >= 500 && attempts < maxRetries {
			log.Printf("Backend %s returned %d, retrying...", backendAddr, status)
			attempts++
			continue
		}

		// Content-Type: text/plain; charset=utf-8
		// Date: Thu, 20 Jun 2025 08:11:30 GMT
		// Content-Length: 24

		// All headers set by the backend (like Content-Type, Set-Cookie, etc.) will be preserved
		// No special retry headers are added by default
		for k, v := range rec.Header() {
			for _, vv := range v {
				w.Header().Add(k, vv)
				
			}
		}
		
		w.Header().Set("X-Retry-Attempts", fmt.Sprintf("%d", attempts))
		w.Header().Set("X-Processed-By", "Go-LB")
		w.Header().Set("X-Backend", backendAddr)
		w.WriteHeader(status)
		rec.Body.WriteTo(w)

		// Get client IP to write your logs (additional feature)
		// var clientIP string
		clientIP := r.RemoteAddr // this is safer
		if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			clientIP = ip
		}

		log.Printf("[%s] %d %s %s → %s (retry=%d) [%s]",
			start.Format("2006-01-02 15:04:05"),
			status,
			r.Method,
			r.URL.Path,
			backendAddr,
			attempts,
			clientIP,
		)
		break
	}
}

// You want to accept an HTTP request from a user, and instead of processing it yourself, you want to forward it to another backend server
// — and then send the backends response back to the user.
// This is called a reverse proxy
// using NewSingleHostReverseProxy, You dont need to manually copy headers, body
// This proxy object will do it for you

// but theres one catch
// Every HTTP request includes a Host header
// GET /login HTTP/1.1
// Host: myapp.com
// When your reverse proxy forwards the request to the backend, by default, it overwrites this Host with:
// Host: localhost:9001
// Your backend might expect the original domain (myapp.com)
// Your backend might use Host-based routing (one backend for api.site.com, another for admin.site.com)
// Logs and analytics would show wrong Host values

// Every httputil.ReverseProxy has something called a Director. Its a function that gets called before each request is forwarded.
// By default, Gos proxy sets:
// req.URL.Host
// req.URL.Scheme
// req.Host = backend address ← this is what we dont want

// proxy.Director = func(req *http.Request) {
// 	originalDirector(req)
// 	req.Host = r.Host
// }

// send it
