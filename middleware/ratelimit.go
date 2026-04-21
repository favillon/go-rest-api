package middleware

import (
	"fmt"
	"sync"
	"time"

	apierrors "backend-productos/api/errors"
	apiresponse "backend-productos/api/response"

	"github.com/gin-gonic/gin"
)

type clientWindow struct {
	count   int
	resetAt time.Time
}

// RateLimitByIP aplica un limite fijo por IP en una ventana de tiempo.
func RateLimitByIP(maxRequests int, window time.Duration) gin.HandlerFunc {
	var (
		mu          sync.Mutex
		clients     = make(map[string]*clientWindow)
		requestTick int
	)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		requestTick++

		if requestTick%200 == 0 {
			for k, v := range clients {
				if now.After(v.resetAt.Add(window)) {
					delete(clients, k)
				}
			}
		}

		entry, ok := clients[ip]
		if !ok || now.After(entry.resetAt) {
			entry = &clientWindow{count: 0, resetAt: now.Add(window)}
			clients[ip] = entry
		}

		remainingBeforeBlock := maxRequests - entry.count
		if remainingBeforeBlock <= 0 {
			retryAfter := int(time.Until(entry.resetAt).Seconds())
			if retryAfter < 1 {
				retryAfter = 1
			}
			mu.Unlock()

			c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
			c.Header("X-RateLimit-Remaining", "0")
			apiresponse.RespondError(
				c,
				429,
				"Demasiadas solicitudes",
				apierrors.RateLimitExceeded,
				"Se excedio el limite de solicitudes. Intenta nuevamente en unos segundos.",
			)
			c.Abort()
			return
		}

		entry.count++
		remaining := maxRequests - entry.count
		if remaining < 0 {
			remaining = 0
		}
		mu.Unlock()

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Next()
	}
}
