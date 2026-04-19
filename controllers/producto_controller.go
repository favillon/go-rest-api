package controllers

import (
	"errors"
	"net/http"

	apierrors "backend-productos/api/errors"
	apiresponse "backend-productos/api/response"
	"backend-productos/config"
	"backend-productos/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Listar todos
func ObtenerProductos(c *gin.Context) {
	if config.DB == nil {
		apiresponse.RespondError(c, http.StatusInternalServerError, "No fue posible procesar la solicitud", apierrors.DBNotInitialized, "database is not initialized")
		return
	}

	var productos []models.Producto
	if err := config.DB.Find(&productos).Error; err != nil {
		apiresponse.RespondError(c, http.StatusInternalServerError, "No fue posible obtener los productos", apierrors.Database, err.Error())
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

		apiresponse.RespondError(c, http.StatusInternalServerError, "No fue posible obtener el producto", apierrors.Database, err.Error())
		return
	}

	apiresponse.RespondSuccess(c, http.StatusOK, "Datos recuperados correctamente", producto)
}

// Crear
func CrearProducto(c *gin.Context) {
	var input models.Producto
	if err := c.ShouldBindJSON(&input); err != nil {
		apiresponse.RespondError(c, http.StatusBadRequest, "La solicitud contiene datos inválidos", apierrors.Validation, err.Error())
		return
	}
	config.DB.Create(&input)
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

		apiresponse.RespondError(c, http.StatusInternalServerError, "No fue posible actualizar el producto", apierrors.Database, err.Error())
		return
	}

	var input models.Producto
	if err := c.ShouldBindJSON(&input); err != nil {
		apiresponse.RespondError(c, http.StatusBadRequest, "La solicitud contiene datos inválidos", apierrors.Validation, err.Error())
		return
	}

	config.DB.Model(&producto).Updates(input)
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

		apiresponse.RespondError(c, http.StatusInternalServerError, "No fue posible eliminar el producto", apierrors.Database, err.Error())
		return
	}

	config.DB.Delete(&producto)
	apiresponse.RespondSuccess(c, http.StatusOK, "Recurso eliminado correctamente", gin.H{"deleted": true})
}
