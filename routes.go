package main

import (
	"backend-productos/controllers"

	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.GET("/productos", controllers.ObtenerProductos)
		api.GET("/productos/:id", controllers.ObtenerProductoPorID)
		api.POST("/productos", controllers.CrearProducto)
		api.PUT("/productos/:id", controllers.ActualizarProducto)
		api.DELETE("/productos/:id", controllers.EliminarProducto)
	}
}
