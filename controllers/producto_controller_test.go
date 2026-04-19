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

	apierrors "backend-productos/api/errors"
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

type apiErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type apiResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    json.RawMessage   `json:"data"`
	Error   *apiErrorResponse `json:"error"`
}

func decodeAPIResponse(t *testing.T, body []byte) apiResponse {
	t.Helper()

	var response apiResponse
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)

	return response
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

	response := decodeAPIResponse(t, w.Body.Bytes())

	var data models.Producto
	err := json.Unmarshal(response.Data, &data)
	assert.NoError(t, err)

	// 3. Assert (Verificar)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Recurso creado correctamente", response.Message)
	assert.Equal(t, "Producto Test", data.Nombre)
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

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, "La solicitud contiene datos inválidos", response.Message)
	assert.Equal(t, string(apierrors.Validation), response.Error.Code)
	assert.Contains(t, response.Error.Message, "Nombre")
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

	response := decodeAPIResponse(t, w.Body.Bytes())

	var data models.Producto
	err := json.Unmarshal(response.Data, &data)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Recurso actualizado correctamente", response.Message)
	assert.Equal(t, uuid.MustParse(productoID), data.ID)
	assert.Equal(t, "Producto actualizado", data.Nombre)
	assert.Equal(t, "Descripcion nueva", data.Descripcion)
	assert.Equal(t, 20.75, data.Precio)
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

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.InvalidParam), response.Error.Code)
	assert.Equal(t, "id inválido", response.Error.Message)
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

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.NotFound), response.Error.Code)
	assert.Equal(t, "producto no encontrado", response.Error.Message)
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

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.Validation), response.Error.Code)
	assert.True(t,
		strings.Contains(response.Error.Message, "invalid character") ||
			strings.Contains(response.Error.Message, "unexpected EOF"),
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

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.Database), response.Error.Code)
	assert.Equal(t, "db failure", response.Error.Message)
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

	response := decodeAPIResponse(t, w.Body.Bytes())

	var data []models.Producto
	err := json.Unmarshal(response.Data, &data)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Datos recuperados correctamente", response.Message)
	assert.Len(t, data, 2)
	assert.Equal(t, "Teclado", data[0].Nombre)
	assert.Equal(t, 99.99, data[0].Precio)
	assert.Equal(t, "Mouse", data[1].Nombre)
	assert.Equal(t, 49.50, data[1].Precio)
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

	response := decodeAPIResponse(t, w.Body.Bytes())

	var data models.Producto
	err := json.Unmarshal(response.Data, &data)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Datos recuperados correctamente", response.Message)
	assert.Equal(t, uuid.MustParse(productoID), data.ID)
	assert.Equal(t, "Monitor", data.Nombre)
	assert.Equal(t, "4K", data.Descripcion)
	assert.Equal(t, 299.90, data.Precio)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestObtenerProductoPorID_IDInvalido(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/productos/:id", ObtenerProductoPorID)

	req, _ := http.NewRequest("GET", "/productos/id-invalido", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.InvalidParam), response.Error.Code)
	assert.Equal(t, "id inválido", response.Error.Message)
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

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.NotFound), response.Error.Code)
	assert.Equal(t, "producto no encontrado", response.Error.Message)
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

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.DBNotInitialized), response.Error.Code)
	assert.Equal(t, "database is not initialized", response.Error.Message)
}

func TestObtenerProductoPorID_DBNoInicializada(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.DB = nil

	r := gin.Default()
	r.GET("/productos/:id", ObtenerProductoPorID)

	req, _ := http.NewRequest("GET", "/productos/550e8400-e29b-41d4-a716-446655440022", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.DBNotInitialized), response.Error.Code)
	assert.Equal(t, "database is not initialized", response.Error.Message)
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

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.Database), response.Error.Code)
	assert.Equal(t, "db failure", response.Error.Message)
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

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.Database), response.Error.Code)
	assert.Equal(t, "db failure", response.Error.Message)
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

	response := decodeAPIResponse(t, w.Body.Bytes())

	var data map[string]bool
	err := json.Unmarshal(response.Data, &data)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Recurso eliminado correctamente", response.Message)
	assert.Equal(t, true, data["deleted"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEliminarProducto_IDInvalido(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.DELETE("/productos/:id", EliminarProducto)

	req, _ := http.NewRequest("DELETE", "/productos/id-invalido", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.InvalidParam), response.Error.Code)
	assert.Equal(t, "id inválido", response.Error.Message)
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

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.NotFound), response.Error.Code)
	assert.Equal(t, "producto no encontrado", response.Error.Message)
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

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.Database), response.Error.Code)
	assert.Equal(t, "db failure", response.Error.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}
