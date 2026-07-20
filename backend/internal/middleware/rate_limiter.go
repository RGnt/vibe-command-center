package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// rateLimitEntry tracks request counts over a rolling window.
type rateLimitEntry struct {
	count     int
	windowEnd time.Time
}

// RateLimiter provides a simple in-memory IP-based rate limiter.
// It is suitable for low-traffic endpoints. For high-throughput deployments
// consider using a distributed store (e.g., Redis).
type RateLimiter struct {
	mu       sync.Mutex
	entries  map[string]*rateLimitEntry
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a RateLimiter allowing `limit` requests per `window` duration.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		entries: make(map[string]*rateLimitEntry),
		limit:   limit,
		window:  window,
	}
	// Periodically clean up stale entries to prevent unbounded memory growth.
	go func() {
		ticker := time.NewTicker(window * 10)
		defer ticker.Stop()
		for range ticker.C {
			rl.mu.Lock()
			now := time.Now()
			for ip, e := range rl.entries {
				if now.After(e.windowEnd) {
					delete(rl.entries, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()
	return rl
}

// Limit is an HTTP middleware that enforces the rate limit.
// Requests from IPs exceeding the limit receive a 429 Too Many Requests response.
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := realIP(r)

		rl.mu.Lock()
		e, ok := rl.entries[ip]
		now := time.Now()
		if !ok || now.After(e.windowEnd) {
			rl.entries[ip] = &rateLimitEntry{count: 1, windowEnd: now.Add(rl.window)}
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}
		e.count++
		if e.count > rl.limit {
			rl.mu.Unlock()
			http.Error(w, "Too many requests — please try again later", http.StatusTooManyRequests)
			return
		}
		rl.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

// realIP extracts the client IP considering potential proxy headers.
// Note: In production, ensure these headers are set by a trusted upstream proxy
// to prevent client spoofing.
func realIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		// Take the first (leftmost) IP — the original client
		parts := strings.SplitN(forwarded, ",", 2)
		return strings.TrimSpace(parts[0])
	}
	// Strip port from RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
