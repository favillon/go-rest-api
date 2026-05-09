package main

import (
	"fmt"
	"log"
	"os"

	"backend-productos/config"
	"backend-productos/graph"
	"backend-productos/internal/application/service"
	"backend-productos/internal/infrastructure/persistence/mongodb"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env.docker"); err != nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("No .env or .env.docker file found, using system environment variables")
		}
	}

	if err := config.InitMongoDB(); err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}

	defer func() {
		if err := config.CloseMongoDB(); err != nil {
			log.Printf("error closing MongoDB: %v", err)
		}
	}()

	productoRepo := mongodb.NewProductoRepository(config.MongoDatabase)
	categoriaRepo := mongodb.NewCategoriaRepository(config.MongoDatabase)
	inventarioRepo := mongodb.NewInventarioRepository(config.MongoDatabase)

	productoSvc := service.NewProductoService(productoRepo)
	categoriaSvc := service.NewCategoriaService(categoriaRepo)
	inventarioSvc := service.NewInventarioService(inventarioRepo)

	resolver := &graph.Resolver{
		ProductoService:   productoSvc,
		CategoriaService:  categoriaSvc,
		InventarioService: inventarioSvc,
	}

	r := gin.Default()

	gqlHandler := handler.NewDefaultServer(
		graph.NewExecutableSchema(graph.Config{Resolvers: resolver}),
	)
	r.POST("/api/v1/graphql", func(c *gin.Context) {
		gqlHandler.ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/api/v1/graphql", gin.WrapH(playground.Handler("GraphQL Playground", "/api/v1/graphql")))

	fmt.Println("backend-productos iniciado")
	fmt.Println("Conexion a MongoDB exitosa")
	fmt.Println("GraphQL Playground: http://localhost:" + getPort() + "/api/v1/graphql")

	if err := r.Run(":" + getPort()); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func getPort() string {
	port := os.Getenv("PORT_APP")
	if port == "" {
		port = "8082"
	}
	return port
}
