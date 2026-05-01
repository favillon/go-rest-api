package main

import (
	"os"
	"strconv"
	"time"

	"backend-productos/internal/infrastructure/http"
	"backend-productos/middleware"

	"github.com/gin-gonic/gin"
)

func parseRateLimit(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 {
			return parsed
		}
	}
	return defaultVal
}

func registerRoutes(r *gin.Engine, handler *http.ProductoHandler) {
	readLimitPerMinute := parseRateLimit("RATE_LIMIT_READ_PER_MINUTE", 30)
	writeLimitPerMinute := parseRateLimit("RATE_LIMIT_WRITE_PER_MINUTE", 10)

	readLimiter := middleware.RateLimitByIP(readLimitPerMinute, time.Minute)
	writeLimiter := middleware.RateLimitByIP(writeLimitPerMinute, time.Minute)

	api := r.Group("/api/v1")
	{
		api.GET("/productos", readLimiter, handler.ObtenerProductos)
		api.GET("/productos/:id", readLimiter, handler.ObtenerProductoPorID)
		api.POST("/productos", writeLimiter, handler.CrearProducto)
		api.PUT("/productos/:id", writeLimiter, handler.ActualizarProducto)
		api.DELETE("/productos/:id", writeLimiter, handler.EliminarProducto)
	}
}
