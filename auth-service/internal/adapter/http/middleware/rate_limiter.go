package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/shared/response"
)

type rateLimiter struct {
	attempts map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		attempts: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Cleanup old entries every minute
	go rl.cleanup()

	return rl
}

func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, times := range rl.attempts {
			// Remove expired attempts
			valid := []time.Time{}
			for _, t := range times {
				if now.Sub(t) < rl.window {
					valid = append(valid, t)
				}
			}
			if len(valid) == 0 {
				delete(rl.attempts, key)
			} else {
				rl.attempts[key] = valid
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use IP address as key (or email from request body for login)
		key := c.ClientIP()

		rl.mu.Lock()
		defer rl.mu.Unlock()

		now := time.Now()

		// Get attempts for this key
		attempts := rl.attempts[key]

		// Remove old attempts outside the window
		validAttempts := []time.Time{}
		for _, t := range attempts {
			if now.Sub(t) < rl.window {
				validAttempts = append(validAttempts, t)
			}
		}

		// Check if limit exceeded
		if len(validAttempts) >= rl.limit {
			response.Error(c, http.StatusTooManyRequests, "TOO_MANY_REQUESTS", "Too many login attempts. Please try again later")
			c.Abort()
			return
		}

		// Add current attempt
		validAttempts = append(validAttempts, now)
		rl.attempts[key] = validAttempts

		c.Next()
	}
}
