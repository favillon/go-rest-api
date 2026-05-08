package model

import (
	"time"

	"github.com/google/uuid"
)

// Categoria represents a product category
type Categoria struct {
	ID          uuid.UUID  `json:"id" bson:"_id"`
	Nombre      string     `json:"nombre" bson:"nombre"`
	Descripcion string     `json:"descripcion" bson:"descripcion"`
	ProductoIDs []string   `json:"producto_ids" bson:"producto_ids"`
	CreatedAt   time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" bson:"updated_at"`
	DeletedAt   *time.Time `json:"-" bson:"deleted_at,omitempty"`
}
