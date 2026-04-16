package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type RateLimitMiddleware struct {
	ipRateLimiter *IPRateLimiter
}

type IPRateLimiter struct {
	limiters map[string]*RateLimiter
	mutex    sync.Mutex
}

type RateLimiter struct {
	tokens         float64   // Current number of tokens
	maxTokens      float64   // Maximum tokens allowed
	refillRate     float64   // Tokens added per second
	lastRefillTime time.Time // Last time tokens were refilled
	lastAccessTime time.Time // Last time this limiter was accessed
	mutex          sync.Mutex
}

func NewRateLimiter(maxTokens, refillRate float64) *RateLimiter {
	now := time.Now()
	return &RateLimiter{
		tokens:         maxTokens,
		maxTokens:      maxTokens,
		refillRate:     refillRate,
		lastRefillTime: now,
		lastAccessTime: now,
	}
}
func NewIPRateLimiter() *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*RateLimiter),
	}
}

func NewRateLimitMiddleware() *RateLimitMiddleware {
	return &RateLimitMiddleware{
		ipRateLimiter: NewIPRateLimiter(),
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *RateLimiter {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	i.cleanupOldLimiters()

	limiter, exists := i.limiters[ip]
	if !exists {
		// Allow 3 requests per minute
		limiter = NewRateLimiter(3, 0.05)
		i.limiters[ip] = limiter
	} else {
		// Update last access time
		limiter.lastAccessTime = time.Now()
	}

	return limiter
}

// cleanupOldLimiters removes entries not accessed in the last 24 hours (call while holding mutex)
func (i *IPRateLimiter) cleanupOldLimiters() {
	cutoff := time.Now().Add(-24 * time.Hour)
	for ip, limiter := range i.limiters {
		if limiter.lastAccessTime.Before(cutoff) {
			delete(i.limiters, ip)
		}
	}
}
func (r *RateLimiter) Allow() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.refillTokens()

	if r.tokens >= 1 {
		r.tokens--
		return true
	}
	return false
}

func (r *RateLimiter) refillTokens() {
	now := time.Now()
	duration := now.Sub(r.lastRefillTime).Seconds()
	tokensToAdd := duration * r.refillRate

	r.tokens += tokensToAdd
	if r.tokens > r.maxTokens {
		r.tokens = r.maxTokens
	}
	r.lastRefillTime = now
}

func (rateLimitMiddleware *RateLimitMiddleware) RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Invalid IP", http.StatusInternalServerError)
			return
		}

		limiter := rateLimitMiddleware.ipRateLimiter.GetLimiter(ip)
		if limiter.Allow() {
			next(w, r)
		} else {
			http.Error(w, "Rate Limit Exceeded", http.StatusTooManyRequests)
		}
	}
}
