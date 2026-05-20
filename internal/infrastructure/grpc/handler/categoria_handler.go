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

// CategoriaHandler implements pb.CategoriaServiceServer.
type CategoriaHandler struct {
	pb.UnimplementedCategoriaServiceServer
	svc *service.CategoriaService
}

// NewCategoriaHandler creates a new gRPC handler for CategoriaService.
func NewCategoriaHandler(svc *service.CategoriaService) *CategoriaHandler {
	return &CategoriaHandler{svc: svc}
}

func toProtoCategoria(c *model.Categoria) *pb.Categoria {
	return &pb.Categoria{
		Id:          c.ID.String(),
		Nombre:      c.Nombre,
		Descripcion: c.Descripcion,
		ProductoIds: c.ProductoIDs,
		CreatedAt:   timestamppb.New(c.CreatedAt),
		UpdatedAt:   timestamppb.New(c.UpdatedAt),
	}
}

// CreateCategoria creates a new category.
func (h *CategoriaHandler) CreateCategoria(ctx context.Context, req *pb.CreateCategoriaRequest) (*pb.CreateCategoriaResponse, error) {
	c := &model.Categoria{
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		ProductoIDs: req.ProductoIds,
	}
	if err := h.svc.Create(ctx, c); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create categoria: %v", err)
	}
	return &pb.CreateCategoriaResponse{Categoria: toProtoCategoria(c)}, nil
}

// GetCategoria retrieves a category by ID.
func (h *CategoriaHandler) GetCategoria(ctx context.Context, req *pb.GetCategoriaRequest) (*pb.GetCategoriaResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid: %v", err)
	}
	c, err := h.svc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "categoria not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get categoria: %v", err)
	}
	return &pb.GetCategoriaResponse{Categoria: toProtoCategoria(c)}, nil
}

// UpdateCategoria updates an existing category.
func (h *CategoriaHandler) UpdateCategoria(ctx context.Context, req *pb.UpdateCategoriaRequest) (*pb.UpdateCategoriaResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid: %v", err)
	}
	c := &model.Categoria{
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		ProductoIDs: req.ProductoIds,
	}
	updated, err := h.svc.Update(ctx, id, c)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "categoria not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update categoria: %v", err)
	}
	return &pb.UpdateCategoriaResponse{Categoria: toProtoCategoria(updated)}, nil
}

// DeleteCategoria performs a soft delete of a category.
func (h *CategoriaHandler) DeleteCategoria(ctx context.Context, req *pb.DeleteCategoriaRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid: %v", err)
	}
	if err := h.svc.Delete(ctx, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "categoria not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete categoria: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// ListCategorias returns all categories.
func (h *CategoriaHandler) ListCategorias(ctx context.Context, _ *pb.ListCategoriasRequest) (*pb.ListCategoriasResponse, error) {
	categorias, err := h.svc.GetAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list categorias: %v", err)
	}
	var protoCategorias []*pb.Categoria
	for _, c := range categorias {
		protoCategorias = append(protoCategorias, toProtoCategoria(&c))
	}
	return &pb.ListCategoriasResponse{Categorias: protoCategorias}, nil
}
