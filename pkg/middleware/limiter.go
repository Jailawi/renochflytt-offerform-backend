package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Limit settings
const (
	reqsPerMinute = 5                // allow 5 requests per minute per IP
	burstSize     = 2                // short bursts allowed
	cleanupAfter  = 10 * time.Minute // cleanup inactive IPs
)

// store limiters per IP
var (
	limiters = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(time.Minute), burstSize)
		limiters[ip] = limiter
	}

	return limiter
}

// simple cleanup goroutine to prevent map from growing forever
func CleanupInactiveLimiters() {
	ticker := time.NewTicker(cleanupAfter)
	go func() {
		for range ticker.C {
			mu.Lock()
			for ip, limiter := range limiters {
				if limiter.Allow() { // checks activity indirectly
					delete(limiters, ip)
				}
			}
			mu.Unlock()
		}
	}()
}
