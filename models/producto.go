package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Producto struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Nombre      string         `json:"nombre" binding:"required" gorm:"not null"`
	Descripcion string         `json:"descripcion" gorm:"not null"`
	Precio      float64        `json:"precio" binding:"required,gt=0" gorm:"type:decimal(10,2);not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"` // El "-" oculta el campo en el JSON
}
