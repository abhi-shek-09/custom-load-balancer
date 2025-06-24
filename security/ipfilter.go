package security

import (
	"github.com/abhi-shek-09/load-balancer/config"
	"net"
	"net/http"
	"strings"
)

// Classless Inter-Domain Routing, is a method of allocating IP addresses and routing internet traffic more efficiently.
// It replaced the older classful addressing system, offering greater flexibility and better routing.
// CIDR notation uses a slash (/) followed by a number to indicate the network prefix length.
// For example, 192.168.1.0/24 means the first 24 bits represent the network, and the remaining 8 bits represent the host

func getClientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		parts := strings.Split(fwd, ",")
		return strings.TrimSpace(parts[0])
	}
	// X-Forwarded-For is a header set by reverse proxies
	// contains a comma-separated list of IPs
	// First IP: original client, rest are intermediate proxies
	// Only trust X-Forwarded-For if r.RemoteAddr is one of your internal proxies

	// if that header isnt there, do the below method
	// r.RemoteAddr contains something like "203.0.113.5:54821", net.SplitHostPort() will extract the IP part
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func isAllowedIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	for _, cidr := range config.AppConfig.AllowedCIDRs {
		_, subnet, err := net.ParseCIDR(cidr)
		if err == nil && subnet.Contains(ip) {
			return true
		}
	}

	return false
}

func IPFilterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		if !isAllowedIP(ip) {
			http.Error(w, "Forbidden: IP not allowed", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
