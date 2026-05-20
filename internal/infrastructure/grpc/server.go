package grpc

import (
	"fmt"
	"log"
	"net"

	"backend-productos/internal/application/service"
	"backend-productos/internal/infrastructure/grpc/handler"
	"backend-productos/internal/infrastructure/grpc/interceptor"
	pb "backend-productos/proto"

	"google.golang.org/grpc"
)

// Server wraps a gRPC server with registered services.
type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

// NewServer creates and configures a gRPC server with all services registered.
func NewServer(
	productoSvc *service.ProductoService,
	categoriaSvc *service.CategoriaService,
	inventarioSvc *service.InventarioService,
) *Server {
	rec := interceptor.RecoveryInterceptor
	log := interceptor.LoggingInterceptor
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(rec, log),
	)

	pb.RegisterProductoServiceServer(s, handler.NewProductoHandler(productoSvc))
	pb.RegisterCategoriaServiceServer(s, handler.NewCategoriaHandler(categoriaSvc))
	pb.RegisterInventarioServiceServer(s, handler.NewInventarioHandler(inventarioSvc))

	return &Server{grpcServer: s}
}

// Start listens on the given port and serves gRPC requests.
func (s *Server) Start(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}
	s.listener = lis

	log.Printf("gRPC server listening on port %s", port)
	return s.grpcServer.Serve(lis)
}

// GracefulStop stops the server gracefully.
func (s *Server) GracefulStop() {
	if s.grpcServer != nil {
		log.Println("shutting down gRPC server gracefully...")
		s.grpcServer.GracefulStop()
	}
}
