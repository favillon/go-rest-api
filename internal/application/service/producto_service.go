package service

import (
	"context"
	"errors"

	"backend-productos/internal/domain"
	"backend-productos/internal/domain/model"
	"backend-productos/internal/domain/port"

	"github.com/google/uuid"
)

type ProductoService struct {
	repo port.ProductoRepository
}

func NewProductoService(repo port.ProductoRepository) *ProductoService {
	return &ProductoService{repo: repo}
}

func (s *ProductoService) GetAll(ctx context.Context, page, limit int) ([]model.Producto, error) {
	return s.repo.GetAll(ctx, page, limit)
}

func (s *ProductoService) GetByID(ctx context.Context, id uuid.UUID) (*model.Producto, error) {
	producto, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, err
		}
		return nil, err
	}
	return producto, nil
}

func (s *ProductoService) Create(ctx context.Context, p *model.Producto) error {
	return s.repo.Create(ctx, p)
}

func (s *ProductoService) Update(ctx context.Context, id uuid.UUID, p *model.Producto) (*model.Producto, error) {
	existente, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existente.Nombre = p.Nombre
	existente.Descripcion = p.Descripcion
	existente.Precio = p.Precio

	if err := s.repo.Update(ctx, existente); err != nil {
		return nil, err
	}

	return existente, nil
}

func (s *ProductoService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}
