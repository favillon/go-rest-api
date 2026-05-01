package main

import (
	"fmt"
	"log"
	"os"

	"backend-productos/config"
	"backend-productos/internal/application/service"
	"backend-productos/internal/infrastructure/http"
	"backend-productos/internal/infrastructure/persistence"

	_ "backend-productos/docs" // Swagger docs auto-generated

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Backend Productos API
// @version 1.0
// @description API REST para gestion de productos con PostgreSQL
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8082
// @basePath /
// @schemes http https
func main() {
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

	repo := persistence.NewProductoRepository(config.DB)
	svc := service.NewProductoService(repo)
	handler := http.NewProductoHandler(svc)

	r := gin.Default()
	registerRoutes(r, handler)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Println("backend-productos iniciado")
	fmt.Println("Conexion a PostgreSQL exitosa")

	port := os.Getenv("PORT_APP")
	if port == "" {
		port = "8082"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
