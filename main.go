package main

import (
	"fmt"
	"log"
	"os"

	"backend-productos/config"
	"backend-productos/controllers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Try to load .env.docker (for Docker environment)
	// then .env (for local development)
	// If neither exists, use system environment variables
	if err := godotenv.Load(".env.docker"); err != nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("No .env or .env.docker file found, using system environment variables")
		}
	}

	if err := config.InitDB(); err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	defer func() {
		if err := config.CloseDB(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/productos", controllers.ObtenerProductos)
		api.GET("/productos/:id", controllers.ObtenerProductoPorID)
		api.POST("/productos", controllers.CrearProducto)
		api.PUT("/productos/:id", controllers.ActualizarProducto)
		api.DELETE("/productos/:id", controllers.EliminarProducto)
	}

	fmt.Println("✓ backend-productos iniciado")
	fmt.Println("✓ Conexión a PostgreSQL exitosa")

	port := os.Getenv("PORT_APP")
	if port == "" {
		port = "8082"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
