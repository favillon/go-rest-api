package controllers

import (
	"errors"
	"log"
	"net/http"

	apierrors "backend-productos/api/errors"
	apiresponse "backend-productos/api/response"
	"backend-productos/config"
	"backend-productos/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	publicDatabaseErrorDetail   = "internal server error"
	publicValidationErrorDetail = "invalid request payload"
)

// Listar todos
func ObtenerProductos(c *gin.Context) {
	if config.DB == nil {
		apiresponse.RespondError(c, http.StatusInternalServerError, "No fue posible procesar la solicitud", apierrors.DBNotInitialized, "database is not initialized")
		return
	}

	var productos []models.Producto
	if err := config.DB.Find(&productos).Error; err != nil {
		log.Printf("database error in ObtenerProductos: %v", err)
		apiresponse.RespondError(c, http.StatusInternalServerError, "No fue posible obtener los productos", apierrors.Database, publicDatabaseErrorDetail)
		return
	}

	apiresponse.RespondSuccess(c, http.StatusOK, "Datos recuperados correctamente", productos)
}

// Obtener por ID
func ObtenerProductoPorID(c *gin.Context) {
	if config.DB == nil {
		apiresponse.RespondError(c, http.StatusInternalServerError, "No fue posible procesar la solicitud", apierrors.DBNotInitialized, "database is not initialized")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.RespondError(c, http.StatusBadRequest, "El parámetro 'id' es obligatorio y debe ser un UUID válido", apierrors.InvalidParam, "id inválido")
		return
	}

	var producto models.Producto
	if err := config.DB.Where("id = ?", id.String()).Take(&producto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			apiresponse.RespondError(c, http.StatusNotFound, "No se encontró el recurso solicitado", apierrors.NotFound, "producto no encontrado")
			return
		}

		log.Printf("database error in ObtenerProductoPorID: %v", err)
		apiresponse.RespondError(c, http.StatusInternalServerError, "No fue posible obtener el producto", apierrors.Database, publicDatabaseErrorDetail)
		return
	}

	apiresponse.RespondSuccess(c, http.StatusOK, "Datos recuperados correctamente", producto)
}

// Crear
func CrearProducto(c *gin.Context) {
	var input models.Producto
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("validation error in CrearProducto: %v", err)
		apiresponse.RespondError(c, http.StatusBadRequest, "La solicitud contiene datos inválidos", apierrors.Validation, publicValidationErrorDetail)
		return
	}
	if err := config.DB.Create(&input).Error; err != nil {
		log.Printf("database error in CrearProducto: %v", err)
		apiresponse.RespondError(c, http.StatusInternalServerError, "No fue posible crear el producto", apierrors.Database, publicDatabaseErrorDetail)
		return
	}
	apiresponse.RespondSuccess(c, http.StatusCreated, "Recurso creado correctamente", input)
}

// Actualizar
func ActualizarProducto(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.RespondError(c, http.StatusBadRequest, "El parámetro 'id' es obligatorio y debe ser un UUID válido", apierrors.InvalidParam, "id inválido")
		return
	}

	var producto models.Producto
	if err := config.DB.Where("id = ?", id.String()).Take(&producto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			apiresponse.RespondError(c, http.StatusNotFound, "No se encontró el recurso solicitado", apierrors.NotFound, "producto no encontrado")
			return
		}

		log.Printf("database error loading producto in ActualizarProducto: %v", err)
		apiresponse.RespondError(c, http.StatusInternalServerError, "No fue posible actualizar el producto", apierrors.Database, publicDatabaseErrorDetail)
		return
	}

	var input models.Producto
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("validation error in ActualizarProducto: %v", err)
		apiresponse.RespondError(c, http.StatusBadRequest, "La solicitud contiene datos inválidos", apierrors.Validation, publicValidationErrorDetail)
		return
	}

	if err := config.DB.Model(&producto).Updates(input).Error; err != nil {
		log.Printf("database error updating producto in ActualizarProducto: %v", err)
		apiresponse.RespondError(c, http.StatusInternalServerError, "No fue posible actualizar el producto", apierrors.Database, publicDatabaseErrorDetail)
		return
	}
	apiresponse.RespondSuccess(c, http.StatusOK, "Recurso actualizado correctamente", producto)
}

// Eliminar
func EliminarProducto(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiresponse.RespondError(c, http.StatusBadRequest, "El parámetro 'id' es obligatorio y debe ser un UUID válido", apierrors.InvalidParam, "id inválido")
		return
	}

	var producto models.Producto
	if err := config.DB.Where("id = ?", id.String()).Take(&producto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			apiresponse.RespondError(c, http.StatusNotFound, "No se encontró el recurso solicitado", apierrors.NotFound, "producto no encontrado")
			return
		}

		log.Printf("database error loading producto in EliminarProducto: %v", err)
		apiresponse.RespondError(c, http.StatusInternalServerError, "No fue posible eliminar el producto", apierrors.Database, publicDatabaseErrorDetail)
		return
	}

	if err := config.DB.Delete(&producto).Error; err != nil {
		log.Printf("database error deleting producto in EliminarProducto: %v", err)
		apiresponse.RespondError(c, http.StatusInternalServerError, "No fue posible eliminar el producto", apierrors.Database, publicDatabaseErrorDetail)
		return
	}
	apiresponse.RespondSuccess(c, http.StatusOK, "Recurso eliminado correctamente", gin.H{"deleted": true})
}
