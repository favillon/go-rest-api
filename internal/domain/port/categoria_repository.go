package port

import (
	"context"

	"backend-productos/internal/domain/model"

	"github.com/google/uuid"
)

// CategoriaRepository defines the contract for category persistence.
type CategoriaRepository interface {
	GetAll(ctx context.Context) ([]model.Categoria, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Categoria, error)
	Create(ctx context.Context, c *model.Categoria) error
	Update(ctx context.Context, c *model.Categoria) error
	Delete(ctx context.Context, id uuid.UUID) error
}
