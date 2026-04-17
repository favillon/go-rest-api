package controllers

import (
	"errors"
	"net/http"

	"backend-productos/config"
	"backend-productos/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Listar todos
func ObtenerProductos(c *gin.Context) {
	if config.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database is not initialized"})
		return
	}

	var productos []models.Producto
	if err := config.DB.Find(&productos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, productos)
}

// Obtener por ID
func ObtenerProductoPorID(c *gin.Context) {
	if config.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database is not initialized"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var producto models.Producto
	if err := config.DB.Where("id = ?", id.String()).Take(&producto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "producto no encontrado"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, producto)
}

// Crear
func CrearProducto(c *gin.Context) {
	var input models.Producto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validacion fallida",
			"message": err.Error(),
		})
		return
	}
	config.DB.Create(&input)
	c.JSON(http.StatusCreated, input)
}

// Actualizar
func ActualizarProducto(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var producto models.Producto
	if err := config.DB.Where("id = ?", id.String()).Take(&producto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "producto no encontrado"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var input models.Producto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&producto).Updates(input)
	c.JSON(http.StatusOK, producto)
}

// Eliminar
func EliminarProducto(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var producto models.Producto
	if err := config.DB.Where("id = ?", id.String()).Take(&producto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "producto no encontrado"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	config.DB.Delete(&producto)
	c.JSON(http.StatusOK, gin.H{"message": "producto eliminado"})
}
