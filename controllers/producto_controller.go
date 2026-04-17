package controllers

import (
	"net/http"

	"backend-productos/config"
	"backend-productos/models"

	"github.com/gin-gonic/gin"
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