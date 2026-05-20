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

// InventarioHandler implements pb.InventarioServiceServer.
type InventarioHandler struct {
	pb.UnimplementedInventarioServiceServer
	svc *service.InventarioService
}

// NewInventarioHandler creates a new gRPC handler for InventarioService.
func NewInventarioHandler(svc *service.InventarioService) *InventarioHandler {
	return &InventarioHandler{svc: svc}
}

func toProtoInventario(i *model.Inventario) *pb.Inventario {
	return &pb.Inventario{
		Id:         i.ID.String(),
		ProductoId: i.ProductoID,
		Cantidad:   int32(i.Cantidad),
		Almacen:    i.Almacen,
		CreatedAt:  timestamppb.New(i.CreatedAt),
		UpdatedAt:  timestamppb.New(i.UpdatedAt),
	}
}

// CreateInventario creates new inventory entry.
func (h *InventarioHandler) CreateInventario(ctx context.Context, req *pb.CreateInventarioRequest) (*pb.CreateInventarioResponse, error) {
	i := &model.Inventario{
		ProductoID: req.ProductoId,
		Cantidad:   int(req.Cantidad),
		Almacen:    req.Almacen,
	}
	if err := h.svc.Create(ctx, i); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create inventario: %v", err)
	}
	return &pb.CreateInventarioResponse{Inventario: toProtoInventario(i)}, nil
}

// GetInventario retrieves inventory by ID.
func (h *InventarioHandler) GetInventario(ctx context.Context, req *pb.GetInventarioRequest) (*pb.GetInventarioResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid: %v", err)
	}
	i, err := h.svc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "inventario not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get inventario: %v", err)
	}
	return &pb.GetInventarioResponse{Inventario: toProtoInventario(i)}, nil
}

// GetInventarioByProductoId retrieves inventory by product ID.
func (h *InventarioHandler) GetInventarioByProductoId(ctx context.Context, req *pb.GetInventarioByProductoIdRequest) (*pb.GetInventarioByProductoIdResponse, error) {
	i, err := h.svc.GetByProductoID(ctx, req.ProductoId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "inventario not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get inventario: %v", err)
	}
	return &pb.GetInventarioByProductoIdResponse{Inventario: toProtoInventario(i)}, nil
}

// UpdateInventario updates an existing inventory entry.
func (h *InventarioHandler) UpdateInventario(ctx context.Context, req *pb.UpdateInventarioRequest) (*pb.UpdateInventarioResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid: %v", err)
	}
	i := &model.Inventario{
		ProductoID: req.ProductoId,
		Cantidad:   int(req.Cantidad),
		Almacen:    req.Almacen,
	}
	updated, err := h.svc.Update(ctx, id, i)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "inventario not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update inventario: %v", err)
	}
	return &pb.UpdateInventarioResponse{Inventario: toProtoInventario(updated)}, nil
}

// DeleteInventario performs a soft delete of an inventory entry.
func (h *InventarioHandler) DeleteInventario(ctx context.Context, req *pb.DeleteInventarioRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid: %v", err)
	}
	if err := h.svc.Delete(ctx, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "inventario not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete inventario: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// ListInventarios returns all inventory entries.
func (h *InventarioHandler) ListInventarios(ctx context.Context, _ *pb.ListInventariosRequest) (*pb.ListInventariosResponse, error) {
	inventarios, err := h.svc.GetAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list inventarios: %v", err)
	}
	var protoInventarios []*pb.Inventario
	for _, i := range inventarios {
		protoInventarios = append(protoInventarios, toProtoInventario(&i))
	}
	return &pb.ListInventariosResponse{Inventarios: protoInventarios}, nil
}
