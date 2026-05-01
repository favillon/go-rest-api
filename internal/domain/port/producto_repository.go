package port

import (
	"context"

	"backend-productos/internal/domain/model"

	"github.com/google/uuid"
)

// ProductoRepository defines the contract for product persistence.
type ProductoRepository interface {
	GetAll(ctx context.Context, page, limit int) ([]model.Producto, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Producto, error)
	Create(ctx context.Context, p *model.Producto) error
	Update(ctx context.Context, p *model.Producto) error
	Delete(ctx context.Context, id uuid.UUID) error
}
