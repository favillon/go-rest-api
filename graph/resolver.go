package graph

import (
	"backend-productos/internal/application/service"
)

type Resolver struct {
	ProductoService   *service.ProductoService
	CategoriaService  *service.CategoriaService
	InventarioService *service.InventarioService
}
