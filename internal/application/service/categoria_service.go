package service

import (
	"context"
	"errors"

	"backend-productos/internal/domain"
	"backend-productos/internal/domain/model"
	"backend-productos/internal/domain/port"

	"github.com/google/uuid"
)

type CategoriaService struct {
	repo port.CategoriaRepository
}

func NewCategoriaService(repo port.CategoriaRepository) *CategoriaService {
	return &CategoriaService{repo: repo}
}

func (s *CategoriaService) GetAll(ctx context.Context) ([]model.Categoria, error) {
	return s.repo.GetAll(ctx)
}

func (s *CategoriaService) GetByID(ctx context.Context, id uuid.UUID) (*model.Categoria, error) {
	categoria, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, err
		}
		return nil, err
	}
	return categoria, nil
}

func (s *CategoriaService) Create(ctx context.Context, c *model.Categoria) error {
	return s.repo.Create(ctx, c)
}

func (s *CategoriaService) Update(ctx context.Context, id uuid.UUID, c *model.Categoria) (*model.Categoria, error) {
	existente, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existente.Nombre = c.Nombre
	existente.Descripcion = c.Descripcion
	existente.ProductoIDs = c.ProductoIDs

	if err := s.repo.Update(ctx, existente); err != nil {
		return nil, err
	}

	return existente, nil
}

func (s *CategoriaService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}
