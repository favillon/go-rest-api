package service

import (
	"context"
	"errors"

	"backend-productos/internal/domain"
	"backend-productos/internal/domain/model"
	"backend-productos/internal/domain/port"

	"github.com/google/uuid"
)

type InventarioService struct {
	repo port.InventarioRepository
}

func NewInventarioService(repo port.InventarioRepository) *InventarioService {
	return &InventarioService{repo: repo}
}

func (s *InventarioService) GetAll(ctx context.Context) ([]model.Inventario, error) {
	return s.repo.GetAll(ctx)
}

func (s *InventarioService) GetByID(ctx context.Context, id uuid.UUID) (*model.Inventario, error) {
	inventario, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, err
		}
		return nil, err
	}
	return inventario, nil
}

func (s *InventarioService) GetByProductoID(ctx context.Context, productoID string) (*model.Inventario, error) {
	return s.repo.GetByProductoID(ctx, productoID)
}

func (s *InventarioService) Create(ctx context.Context, i *model.Inventario) error {
	return s.repo.Create(ctx, i)
}

func (s *InventarioService) Update(ctx context.Context, id uuid.UUID, i *model.Inventario) (*model.Inventario, error) {
	existente, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existente.ProductoID = i.ProductoID
	existente.Cantidad = i.Cantidad
	existente.Almacen = i.Almacen

	if err := s.repo.Update(ctx, existente); err != nil {
		return nil, err
	}

	return existente, nil
}

func (s *InventarioService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}
