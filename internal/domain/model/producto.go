package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Producto represents a product in the system
// @Description Product model with pricing and timestamps
type Producto struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	Nombre      string         `json:"nombre" binding:"required" gorm:"not null" example:"Laptop Dell XPS 13"`
	Descripcion string         `json:"descripcion" gorm:"not null" example:"High-performance ultrabook"`
	Precio      float64        `json:"precio" binding:"required,gt=0" gorm:"type:decimal(10,2);not null" example:"999.99"`
	CreatedAt   time.Time      `json:"created_at" example:"2026-04-20T10:00:00Z"`
	UpdatedAt   time.Time      `json:"updated_at" example:"2026-04-20T10:00:00Z"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
