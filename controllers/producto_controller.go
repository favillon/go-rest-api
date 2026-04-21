package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

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
	defaultPage                 = 1
	defaultLimit                = 20
	maxLimit                    = 100
)

func parsePaginationParams(c *gin.Context) (int, int, error) {
	page := defaultPage
	limit := defaultLimit

	if pageQuery := c.Query("page"); pageQuery != "" {
		parsedPage, err := strconv.Atoi(pageQuery)
		if err != nil || parsedPage < 1 {
			return 0, 0, errors.New("page debe ser un entero mayor o igual a 1")
		}
		page = parsedPage
	}

	if limitQuery := c.Query("limit"); limitQuery != "" {
		parsedLimit, err := strconv.Atoi(limitQuery)
		if err != nil || parsedLimit < 1 {
			return 0, 0, errors.New("limit debe ser un entero mayor o igual a 1")
		}
		if parsedLimit > maxLimit {
			parsedLimit = maxLimit
		}
		limit = parsedLimit
	}

	return page, limit, nil
}

// ObtenerProductos retrieves all products with pagination
// @Summary Get all products
// @Description Retrieve a paginated list of all products. Default page is 1 with 20 items per page (max 100).
// @Tags Productos
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)" default(1)
// @Param limit query int false "Items per page (default 20, max 100)" default(20)
// @Success 200 {object} map[string]interface{} "status:success,message:Datos recuperados correctamente,data:[]"
// @Failure 400 {object} map[string]interface{} "status:error (invalid pagination params)"
// @Failure 500 {object} map[string]interface{} "status:error (database error)"
// @Router /api/v1/productos [get]
func ObtenerProductos(c *gin.Context) {
	if config.DB == nil {
		apiresponse.RespondError(c, http.StatusInternalServerError, "No fue posible procesar la solicitud", apierrors.DBNotInitialized, "database is not initialized")
		return
	}

	page, limit, err := parsePaginationParams(c)
	if err != nil {
		apiresponse.RespondError(c, http.StatusBadRequest, "Parametros de paginacion invalidos", apierrors.InvalidParam, err.Error())
		return
	}

	offset := (page - 1) * limit

	var productos []models.Producto
	if err := config.DB.Limit(limit).Offset(offset).Find(&productos).Error; err != nil {
		log.Printf("database error in ObtenerProductos: %v", err)
		apiresponse.RespondError(c, http.StatusInternalServerError, "No fue posible obtener los productos", apierrors.Database, publicDatabaseErrorDetail)
		return
	}

	apiresponse.RespondSuccess(c, http.StatusOK, "Datos recuperados correctamente", productos)
}

// ObtenerProductoPorID retrieves a single product by ID
// @Summary Get product by ID
// @Description Retrieve a specific product using its UUID
// @Tags Productos
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Success 200 {object} map[string]interface{} "status:success,data:Producto"
// @Failure 400 {object} map[string]interface{} "status:error (invalid UUID)"
// @Failure 404 {object} map[string]interface{} "status:error (product not found)"
// @Failure 500 {object} map[string]interface{} "status:error (database error)"
// @Router /api/v1/productos/{id} [get]
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

// CrearProducto creates a new product
// @Summary Create a new product
// @Description Create a new product with required fields: nombre and precio (must be > 0)
// @Tags Productos
// @Accept json
// @Produce json
// @Param product body models.Producto true "Product data"
// @Success 201 {object} map[string]interface{} "status:success,message:Recurso creado correctamente,data:Producto"
// @Failure 400 {object} map[string]interface{} "status:error (validation error)"
// @Failure 429 {object} map[string]interface{} "status:error (rate limit exceeded)"
// @Failure 500 {object} map[string]interface{} "status:error (database error)"
// @Router /api/v1/productos [post]
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

// ActualizarProducto updates an existing product
// @Summary Update a product
// @Description Update a product by ID. Provide fields to update.
// @Tags Productos
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Param product body models.Producto true "Updated product data"
// @Success 200 {object} map[string]interface{} "status:success,message:Recurso actualizado correctamente,data:Producto"
// @Failure 400 {object} map[string]interface{} "status:error (invalid UUID or validation error)"
// @Failure 404 {object} map[string]interface{} "status:error (product not found)"
// @Failure 429 {object} map[string]interface{} "status:error (rate limit exceeded)"
// @Failure 500 {object} map[string]interface{} "status:error (database error)"
// @Router /api/v1/productos/{id} [put]
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

// EliminarProducto deletes (soft delete) a product
// @Summary Delete a product
// @Description Soft delete a product by ID (marks as deleted, does not remove from database)
// @Tags Productos
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Success 200 {object} map[string]interface{} "status:success,message:Recurso eliminado correctamente,data:deleted:true"
// @Failure 400 {object} map[string]interface{} "status:error (invalid UUID)"
// @Failure 404 {object} map[string]interface{} "status:error (product not found)"
// @Failure 429 {object} map[string]interface{} "status:error (rate limit exceeded)"
// @Failure 500 {object} map[string]interface{} "status:error (database error)"
// @Router /api/v1/productos/{id} [delete]
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
