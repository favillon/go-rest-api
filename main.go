package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"backend-productos/config"
	"backend-productos/internal/application/service"
	grpcserver "backend-productos/internal/infrastructure/grpc"
	"backend-productos/internal/infrastructure/persistence/mongodb"

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

	server := grpcserver.NewServer(productoSvc, categoriaSvc, inventarioSvc)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		server.GracefulStop()
	}()

	port := getGRPCPort()
	fmt.Println("backend-productos iniciado")
	fmt.Println("Conexion a MongoDB exitosa")
	fmt.Printf("gRPC server listening on port %s\n", port)

	if err := server.Start(port); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}

func getGRPCPort() string {
	port := os.Getenv("PORT_GRPC")
	if port == "" {
		port = "50051"
	}
	return port
}
