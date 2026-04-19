package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"backend-productos/config"
	"backend-productos/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupTestDB crea una conexión GORM mockeada para no tocar la DB real
func SetupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn: dbMock,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return db, mock
}

func TestCrearProducto_Exito(t *testing.T) {
	// 1. Arrange (Preparar)
	gin.SetMode(gin.TestMode)
	db, mock := SetupTestDB(t)
	config.DB = db // Inyectamos el mock en la config global

	r := gin.Default()
	r.POST("/productos", CrearProducto)

	productoReq := models.Producto{
		Nombre: "Producto Test",
		Precio: 10.50,
	}
	jsonValue, _ := json.Marshal(productoReq)

	// Esperamos que GORM intente insertar en la tabla productos
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("550e8400-e29b-41d4-a716-446655440000"))
	mock.ExpectCommit()

	// 2. Act (Ejecutar)
	req, _ := http.NewRequest("POST", "/productos", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 3. Assert (Verificar)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCrearProducto_ValidacionFallida(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/productos", CrearProducto)

	// Payload sin nombre (que es requerido)
	payload := []byte(`{"precio": 10.50}`)

	req, _ := http.NewRequest("POST", "/productos", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Validacion fallida", response["error"])
	assert.Contains(t, response["message"], "Nombre")
}

func TestActualizarProducto_Exito(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := SetupTestDB(t)
	config.DB = db

	r := gin.Default()
	r.PUT("/productos/:id", ActualizarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440000"
	now := time.Now()

	// 1) Se consulta el producto existente
	mock.ExpectQuery(`SELECT .* FROM "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "descripcion", "precio", "created_at", "updated_at", "deleted_at"}).
			AddRow(productoID, "Producto viejo", "Descripcion vieja", 10.50, now, now, nil))

	// 2) Se actualiza el producto
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "productos"`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	payload := []byte(`{"nombre":"Producto actualizado","descripcion":"Descripcion nueva","precio":20.75}`)
	req, _ := http.NewRequest("PUT", "/productos/"+productoID, bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response models.Producto
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, uuid.MustParse(productoID), response.ID)
	assert.Equal(t, "Producto actualizado", response.Nombre)
	assert.Equal(t, "Descripcion nueva", response.Descripcion)
	assert.Equal(t, 20.75, response.Precio)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActualizarProducto_IDInvalido(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PUT("/productos/:id", ActualizarProducto)

	payload := []byte(`{"nombre":"Producto actualizado","precio":20.75}`)
	req, _ := http.NewRequest("PUT", "/productos/id-invalido", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "id inválido", response["error"])
}

func TestActualizarProducto_NoEncontrado(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := SetupTestDB(t)
	config.DB = db

	r := gin.Default()
	r.PUT("/productos/:id", ActualizarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440001"

	mock.ExpectQuery(`SELECT .* FROM "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "descripcion", "precio", "created_at", "updated_at", "deleted_at"}))

	payload := []byte(`{"nombre":"Producto actualizado","descripcion":"Descripcion nueva","precio":20.75}`)
	req, _ := http.NewRequest("PUT", "/productos/"+productoID, bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "producto no encontrado", response["error"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActualizarProducto_PayloadInvalido(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := SetupTestDB(t)
	config.DB = db

	r := gin.Default()
	r.PUT("/productos/:id", ActualizarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440002"
	now := time.Now()

	mock.ExpectQuery(`SELECT .* FROM "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "descripcion", "precio", "created_at", "updated_at", "deleted_at"}).
			AddRow(productoID, "Producto viejo", "Descripcion vieja", 10.50, now, now, nil))

	// JSON mal formado para forzar error de binding
	payload := []byte(`{"nombre":"Producto actualizado",`)
	req, _ := http.NewRequest("PUT", "/productos/"+productoID, bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.True(t,
		strings.Contains(response["error"], "invalid character") ||
			strings.Contains(response["error"], "unexpected EOF"),
	)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActualizarProducto_ErrorDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := SetupTestDB(t)
	config.DB = db

	r := gin.Default()
	r.PUT("/productos/:id", ActualizarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440003"
	mock.ExpectQuery(`SELECT .* FROM "productos"`).WillReturnError(errors.New("db failure"))

	payload := []byte(`{"nombre":"Producto actualizado","descripcion":"Descripcion nueva","precio":20.75}`)
	req, _ := http.NewRequest("PUT", "/productos/"+productoID, bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "db failure", response["error"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestObtenerProductos_Exito(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := SetupTestDB(t)
	config.DB = db

	r := gin.Default()
	r.GET("/productos", ObtenerProductos)

	now := time.Now()
	mock.ExpectQuery(`SELECT .* FROM "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "descripcion", "precio", "created_at", "updated_at", "deleted_at"}).
			AddRow("550e8400-e29b-41d4-a716-446655440010", "Teclado", "Mecanico", 99.99, now, now, nil).
			AddRow("550e8400-e29b-41d4-a716-446655440011", "Mouse", "Inalambrico", 49.50, now, now, nil))

	req, _ := http.NewRequest("GET", "/productos", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response []models.Producto
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, response, 2)
	assert.Equal(t, "Teclado", response[0].Nombre)
	assert.Equal(t, 99.99, response[0].Precio)
	assert.Equal(t, "Mouse", response[1].Nombre)
	assert.Equal(t, 49.50, response[1].Precio)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestObtenerProductoPorID_Exito(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := SetupTestDB(t)
	config.DB = db

	r := gin.Default()
	r.GET("/productos/:id", ObtenerProductoPorID)

	productoID := "550e8400-e29b-41d4-a716-446655440020"
	now := time.Now()

	mock.ExpectQuery(`SELECT .* FROM "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "descripcion", "precio", "created_at", "updated_at", "deleted_at"}).
			AddRow(productoID, "Monitor", "4K", 299.90, now, now, nil))

	req, _ := http.NewRequest("GET", "/productos/"+productoID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response models.Producto
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, uuid.MustParse(productoID), response.ID)
	assert.Equal(t, "Monitor", response.Nombre)
	assert.Equal(t, "4K", response.Descripcion)
	assert.Equal(t, 299.90, response.Precio)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestObtenerProductoPorID_IDInvalido(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/productos/:id", ObtenerProductoPorID)

	req, _ := http.NewRequest("GET", "/productos/id-invalido", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "id inválido", response["error"])
}

func TestObtenerProductoPorID_NoEncontrado(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := SetupTestDB(t)
	config.DB = db

	r := gin.Default()
	r.GET("/productos/:id", ObtenerProductoPorID)

	productoID := "550e8400-e29b-41d4-a716-446655440021"

	mock.ExpectQuery(`SELECT .* FROM "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "descripcion", "precio", "created_at", "updated_at", "deleted_at"}))

	req, _ := http.NewRequest("GET", "/productos/"+productoID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "producto no encontrado", response["error"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestObtenerProductos_DBNoInicializada(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.DB = nil

	r := gin.Default()
	r.GET("/productos", ObtenerProductos)

	req, _ := http.NewRequest("GET", "/productos", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "database is not initialized", response["error"])
}

func TestObtenerProductoPorID_DBNoInicializada(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.DB = nil

	r := gin.Default()
	r.GET("/productos/:id", ObtenerProductoPorID)

	req, _ := http.NewRequest("GET", "/productos/550e8400-e29b-41d4-a716-446655440022", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "database is not initialized", response["error"])
}

func TestObtenerProductos_ErrorDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := SetupTestDB(t)
	config.DB = db

	r := gin.Default()
	r.GET("/productos", ObtenerProductos)

	mock.ExpectQuery(`SELECT .* FROM "productos"`).WillReturnError(errors.New("db failure"))

	req, _ := http.NewRequest("GET", "/productos", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "db failure", response["error"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestObtenerProductoPorID_ErrorDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := SetupTestDB(t)
	config.DB = db

	r := gin.Default()
	r.GET("/productos/:id", ObtenerProductoPorID)

	productoID := "550e8400-e29b-41d4-a716-446655440022"
	mock.ExpectQuery(`SELECT .* FROM "productos"`).WillReturnError(errors.New("db failure"))

	req, _ := http.NewRequest("GET", "/productos/"+productoID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "db failure", response["error"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEliminarProducto_Exito(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := SetupTestDB(t)
	config.DB = db

	r := gin.Default()
	r.DELETE("/productos/:id", EliminarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440030"
	now := time.Now()

	mock.ExpectQuery(`SELECT .* FROM "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "descripcion", "precio", "created_at", "updated_at", "deleted_at"}).
			AddRow(productoID, "Camara", "HD", 149.99, now, now, nil))

	// Por soft delete, GORM ejecuta UPDATE sobre deleted_at.
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "productos"`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	req, _ := http.NewRequest("DELETE", "/productos/"+productoID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "producto eliminado", response["message"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEliminarProducto_IDInvalido(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.DELETE("/productos/:id", EliminarProducto)

	req, _ := http.NewRequest("DELETE", "/productos/id-invalido", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "id inválido", response["error"])
}

func TestEliminarProducto_NoEncontrado(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := SetupTestDB(t)
	config.DB = db

	r := gin.Default()
	r.DELETE("/productos/:id", EliminarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440031"
	mock.ExpectQuery(`SELECT .* FROM "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "descripcion", "precio", "created_at", "updated_at", "deleted_at"}))

	req, _ := http.NewRequest("DELETE", "/productos/"+productoID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "producto no encontrado", response["error"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEliminarProducto_ErrorDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := SetupTestDB(t)
	config.DB = db

	r := gin.Default()
	r.DELETE("/productos/:id", EliminarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440032"
	mock.ExpectQuery(`SELECT .* FROM "productos"`).WillReturnError(errors.New("db failure"))

	req, _ := http.NewRequest("DELETE", "/productos/"+productoID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "db failure", response["error"])
	assert.NoError(t, mock.ExpectationsWereMet())
}
