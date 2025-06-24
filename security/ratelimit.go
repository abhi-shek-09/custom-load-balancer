package security

import (
	"net/http"
	"sync"
	"time"
)

type tokenBucket struct {
	tokens         int
	lastRefillTime time.Time
	mu             sync.Mutex
}

var (
	limiters     sync.Map // map[string]*tokenBucket
	maxTokens    = 5
	refillPeriod = 10 * time.Second
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		limiterIface, _ := limiters.LoadOrStore(ip, &tokenBucket{
			tokens:         maxTokens,
			lastRefillTime: time.Now(),
		})
		limiter := limiterIface.(*tokenBucket)

		limiter.mu.Lock()
		defer limiter.mu.Unlock()

		now := time.Now()
		elapsed := now.Sub(limiter.lastRefillTime)
		if elapsed > refillPeriod {
			limiter.tokens = maxTokens
			limiter.lastRefillTime = now
		}

		if limiter.tokens <= 0 {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		limiter.tokens--
		next.ServeHTTP(w, r)
	})
}