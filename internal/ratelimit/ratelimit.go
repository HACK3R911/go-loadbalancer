package ratelimit

import (
	"net"
	"net/http"
	"sync"
)

type Config struct {
	Capacity   float64
	RefillRate float64
}

type RateLimiter struct {
	buckets sync.Map
	config  Config
}

func New(cfg Config) *RateLimiter {
	return &RateLimiter{
		config: cfg,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := extractIP(r)
		bucket, _ := rl.buckets.LoadOrStore(clientIP, newTokenBucket(
			rl.config.Capacity,
			rl.config.RefillRate,
		))

		if !bucket.(*TokenBucket).Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func extractIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
