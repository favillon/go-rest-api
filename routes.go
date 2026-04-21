package main

import (
	"backend-productos/controllers"
	"backend-productos/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	readLimitPerMinute  = 30
	writeLimitPerMinute = 10
)

func registerRoutes(r *gin.Engine) {
	readLimiter := middleware.RateLimitByIP(readLimitPerMinute, time.Minute)
	writeLimiter := middleware.RateLimitByIP(writeLimitPerMinute, time.Minute)

	api := r.Group("/api/v1")
	{
		api.GET("/productos", readLimiter, controllers.ObtenerProductos)
		api.GET("/productos/:id", readLimiter, controllers.ObtenerProductoPorID)
		api.POST("/productos", writeLimiter, controllers.CrearProducto)
		api.PUT("/productos/:id", writeLimiter, controllers.ActualizarProducto)
		api.DELETE("/productos/:id", writeLimiter, controllers.EliminarProducto)
	}
}
