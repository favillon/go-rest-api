package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Producto struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Nombre      string         `json:"nombre" gorm:"not null"`
	Descripcion string         `json:"descripcion"`
	Precio      float64        `json:"precio" gorm:"type:decimal(10,2);not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"` // El "-" oculta el campo en el JSON
}
