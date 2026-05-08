package port

import (
	"context"

	"backend-productos/internal/domain/model"

	"github.com/google/uuid"
)

// InventarioRepository defines the contract for inventory persistence.
type InventarioRepository interface {
	GetAll(ctx context.Context) ([]model.Inventario, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Inventario, error)
	GetByProductoID(ctx context.Context, productoID string) (*model.Inventario, error)
	Create(ctx context.Context, i *model.Inventario) error
	Update(ctx context.Context, i *model.Inventario) error
	Delete(ctx context.Context, id uuid.UUID) error
}
