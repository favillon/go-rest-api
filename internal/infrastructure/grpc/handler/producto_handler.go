package handler

import (
	"context"
	"errors"

	"backend-productos/internal/application/service"
	"backend-productos/internal/domain"
	"backend-productos/internal/domain/model"
	pb "backend-productos/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProductoHandler implements pb.ProductoServiceServer.
type ProductoHandler struct {
	pb.UnimplementedProductoServiceServer
	svc *service.ProductoService
}

// NewProductoHandler creates a new gRPC handler for ProductoService.
func NewProductoHandler(svc *service.ProductoService) *ProductoHandler {
	return &ProductoHandler{svc: svc}
}

func toProtoProducto(p *model.Producto) *pb.Producto {
	return &pb.Producto{
		Id:          p.ID.String(),
		Nombre:      p.Nombre,
		Descripcion: p.Descripcion,
		Precio:      p.Precio,
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}
}

// CreateProducto creates a new product.
func (h *ProductoHandler) CreateProducto(ctx context.Context, req *pb.CreateProductoRequest) (*pb.CreateProductoResponse, error) {
	p := &model.Producto{
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		Precio:      req.Precio,
	}
	if err := h.svc.Create(ctx, p); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create producto: %v", err)
	}
	return &pb.CreateProductoResponse{Producto: toProtoProducto(p)}, nil
}

// GetProducto retrieves a product by ID.
func (h *ProductoHandler) GetProducto(ctx context.Context, req *pb.GetProductoRequest) (*pb.GetProductoResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid: %v", err)
	}
	p, err := h.svc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "producto not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get producto: %v", err)
	}
	return &pb.GetProductoResponse{Producto: toProtoProducto(p)}, nil
}

// UpdateProducto updates an existing product.
func (h *ProductoHandler) UpdateProducto(ctx context.Context, req *pb.UpdateProductoRequest) (*pb.UpdateProductoResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid: %v", err)
	}
	p := &model.Producto{
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		Precio:      req.Precio,
	}
	updated, err := h.svc.Update(ctx, id, p)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "producto not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update producto: %v", err)
	}
	return &pb.UpdateProductoResponse{Producto: toProtoProducto(updated)}, nil
}

// DeleteProducto performs a soft delete of a product.
func (h *ProductoHandler) DeleteProducto(ctx context.Context, req *pb.DeleteProductoRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid: %v", err)
	}
	if err := h.svc.Delete(ctx, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "producto not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete producto: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// ListProductos returns a paginated list of products.
func (h *ProductoHandler) ListProductos(ctx context.Context, req *pb.ListProductosRequest) (*pb.ListProductosResponse, error) {
	page := int(req.Page)
	limit := int(req.Limit)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	productos, err := h.svc.GetAll(ctx, page, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list productos: %v", err)
	}
	var protoProductos []*pb.Producto
	for _, p := range productos {
		protoProductos = append(protoProductos, toProtoProducto(&p))
	}
	return &pb.ListProductosResponse{Productos: protoProductos}, nil
}
