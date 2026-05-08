package model

import (
	"time"

	"github.com/google/uuid"
)

// Inventario represents stock information for a product
type Inventario struct {
	ID         uuid.UUID  `json:"id" bson:"_id"`
	ProductoID string     `json:"producto_id" bson:"producto_id"`
	Cantidad   int        `json:"cantidad" bson:"cantidad"`
	Almacen    string     `json:"almacen" bson:"almacen"`
	CreatedAt  time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" bson:"updated_at"`
	DeletedAt  *time.Time `json:"-" bson:"deleted_at,omitempty"`
}
