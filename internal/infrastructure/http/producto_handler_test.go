package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apierrors "backend-productos/api/errors"
	"backend-productos/internal/application/service"
	"backend-productos/internal/domain/model"
	"backend-productos/internal/infrastructure/persistence"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn: dbMock,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	t.Cleanup(func() {
		_ = dbMock.Close()
	})

	return db, mock
}

func newTestHandler(db *gorm.DB) *ProductoHandler {
	repo := persistence.NewProductoRepository(db)
	svc := service.NewProductoService(repo)
	return NewProductoHandler(svc)
}

func TestCrearProducto_Exito(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.POST("/productos", handler.CrearProducto)

	productoReq := model.Producto{
		Nombre: "Producto Test",
		Precio: 10.50,
	}
	jsonValue, _ := json.Marshal(productoReq)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("550e8400-e29b-41d4-a716-446655440000"))
	mock.ExpectCommit()

	req, _ := http.NewRequest("POST", "/productos", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	var data model.Producto
	err := json.Unmarshal(response.Data, &data)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Recurso creado correctamente", response.Message)
	assert.Equal(t, "Producto Test", data.Nombre)
}

func TestCrearProducto_ValidacionFallida(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.POST("/productos", handler.CrearProducto)

	payload := []byte(`{"precio": 10.50}`)

	req, _ := http.NewRequest("POST", "/productos", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, "La solicitud contiene datos invalidos", response.Message)
	assert.Equal(t, string(apierrors.Validation), response.Error.Code)
	assert.Equal(t, "invalid request payload", response.Error.Message)
}

func TestCrearProducto_ErrorDBCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.POST("/productos", handler.CrearProducto)

	productoReq := model.Producto{
		Nombre: "Producto Test",
		Precio: 10.50,
	}
	jsonValue, _ := json.Marshal(productoReq)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "productos"`).WillReturnError(errors.New("db failure"))
	mock.ExpectRollback()

	req, _ := http.NewRequest("POST", "/productos", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.Database), response.Error.Code)
	assert.Equal(t, "internal server error", response.Error.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActualizarProducto_Exito(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.PUT("/productos/:id", handler.ActualizarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440000"
	now := time.Now()

	mock.ExpectQuery(`SELECT .* FROM "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "descripcion", "precio", "created_at", "updated_at", "deleted_at"}).
			AddRow(productoID, "Producto viejo", "Descripcion vieja", 10.50, now, now, nil))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "productos"`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	payload := []byte(`{"nombre":"Producto actualizado","descripcion":"Descripcion nueva","precio":20.75}`)
	req, _ := http.NewRequest("PUT", "/productos/"+productoID, bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	var data model.Producto
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
	db, _ := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.PUT("/productos/:id", handler.ActualizarProducto)

	payload := []byte(`{"nombre":"Producto actualizado","precio":20.75}`)
	req, _ := http.NewRequest("PUT", "/productos/id-invalido", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.InvalidParam), response.Error.Code)
	assert.Equal(t, "id invalido", response.Error.Message)
}

func TestActualizarProducto_NoEncontrado(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.PUT("/productos/:id", handler.ActualizarProducto)

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
	db, _ := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.PUT("/productos/:id", handler.ActualizarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440002"

	payload := []byte(`{"nombre":"Producto actualizado",`)
	req, _ := http.NewRequest("PUT", "/productos/"+productoID, bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.Validation), response.Error.Code)
	assert.Equal(t, "invalid request payload", response.Error.Message)
}

func TestActualizarProducto_ErrorDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.PUT("/productos/:id", handler.ActualizarProducto)

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
	assert.Equal(t, "internal server error", response.Error.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActualizarProducto_ErrorDBUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.PUT("/productos/:id", handler.ActualizarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440004"
	now := time.Now()

	mock.ExpectQuery(`SELECT .* FROM "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "descripcion", "precio", "created_at", "updated_at", "deleted_at"}).
			AddRow(productoID, "Producto viejo", "Descripcion vieja", 10.50, now, now, nil))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "productos"`).WillReturnError(errors.New("db failure"))
	mock.ExpectRollback()

	payload := []byte(`{"nombre":"Producto actualizado","descripcion":"Descripcion nueva","precio":20.75}`)
	req, _ := http.NewRequest("PUT", "/productos/"+productoID, bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.Database), response.Error.Code)
	assert.Equal(t, "internal server error", response.Error.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestObtenerProductos_Exito(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.GET("/productos", handler.ObtenerProductos)

	now := time.Now()
	mock.ExpectQuery(`SELECT .* FROM "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "descripcion", "precio", "created_at", "updated_at", "deleted_at"}).
			AddRow("550e8400-e29b-41d4-a716-446655440010", "Teclado", "Mecanico", 99.99, now, now, nil).
			AddRow("550e8400-e29b-41d4-a716-446655440011", "Mouse", "Inalambrico", 49.50, now, now, nil))

	req, _ := http.NewRequest("GET", "/productos", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	var data []model.Producto
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

func TestObtenerProductos_PaginacionConParametros(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.GET("/productos", handler.ObtenerProductos)

	now := time.Now()
	mock.ExpectQuery(`SELECT .* FROM "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "descripcion", "precio", "created_at", "updated_at", "deleted_at"}).
			AddRow("550e8400-e29b-41d4-a716-446655440012", "Audifonos", "Noise cancelling", 120.0, now, now, nil))

	req, _ := http.NewRequest("GET", "/productos?page=2&limit=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	var data []model.Producto
	err := json.Unmarshal(response.Data, &data)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Datos recuperados correctamente", response.Message)
	assert.Len(t, data, 1)
	assert.Equal(t, "Audifonos", data[0].Nombre)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestObtenerProductos_PaginacionInvalida(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.GET("/productos", handler.ObtenerProductos)

	req, _ := http.NewRequest("GET", "/productos?page=0&limit=-5", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.InvalidParam), response.Error.Code)
}

func TestObtenerProductoPorID_Exito(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.GET("/productos/:id", handler.ObtenerProductoPorID)

	productoID := "550e8400-e29b-41d4-a716-446655440020"
	now := time.Now()

	mock.ExpectQuery(`SELECT .* FROM "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "descripcion", "precio", "created_at", "updated_at", "deleted_at"}).
			AddRow(productoID, "Monitor", "4K", 299.90, now, now, nil))

	req, _ := http.NewRequest("GET", "/productos/"+productoID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	var data model.Producto
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
	db, _ := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.GET("/productos/:id", handler.ObtenerProductoPorID)

	req, _ := http.NewRequest("GET", "/productos/id-invalido", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.InvalidParam), response.Error.Code)
	assert.Equal(t, "id invalido", response.Error.Message)
}

func TestObtenerProductoPorID_NoEncontrado(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.GET("/productos/:id", handler.ObtenerProductoPorID)

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

func TestObtenerProductos_ErrorDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.GET("/productos", handler.ObtenerProductos)

	mock.ExpectQuery(`SELECT .* FROM "productos"`).WillReturnError(errors.New("db failure"))

	req, _ := http.NewRequest("GET", "/productos", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.Database), response.Error.Code)
	assert.Equal(t, "internal server error", response.Error.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestObtenerProductoPorID_ErrorDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.GET("/productos/:id", handler.ObtenerProductoPorID)

	productoID := "550e8400-e29b-41d4-a716-446655440022"
	mock.ExpectQuery(`SELECT .* FROM "productos"`).WillReturnError(errors.New("db failure"))

	req, _ := http.NewRequest("GET", "/productos/"+productoID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.Database), response.Error.Code)
	assert.Equal(t, "internal server error", response.Error.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEliminarProducto_Exito(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.DELETE("/productos/:id", handler.EliminarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440030"
	now := time.Now()

	mock.ExpectQuery(`SELECT .* FROM "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "descripcion", "precio", "created_at", "updated_at", "deleted_at"}).
			AddRow(productoID, "Camara", "HD", 149.99, now, now, nil))

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
	db, _ := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.DELETE("/productos/:id", handler.EliminarProducto)

	req, _ := http.NewRequest("DELETE", "/productos/id-invalido", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.InvalidParam), response.Error.Code)
	assert.Equal(t, "id invalido", response.Error.Message)
}

func TestEliminarProducto_NoEncontrado(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.DELETE("/productos/:id", handler.EliminarProducto)

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
	db, mock := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.DELETE("/productos/:id", handler.EliminarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440032"
	mock.ExpectQuery(`SELECT .* FROM "productos"`).WillReturnError(errors.New("db failure"))

	req, _ := http.NewRequest("DELETE", "/productos/"+productoID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.Database), response.Error.Code)
	assert.Equal(t, "internal server error", response.Error.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEliminarProducto_ErrorDBDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := setupTestDB(t)
	handler := newTestHandler(db)

	r := gin.Default()
	r.DELETE("/productos/:id", handler.EliminarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440033"
	now := time.Now()

	mock.ExpectQuery(`SELECT .* FROM "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "descripcion", "precio", "created_at", "updated_at", "deleted_at"}).
			AddRow(productoID, "Camara", "HD", 149.99, now, now, nil))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "productos"`).WillReturnError(errors.New("db failure"))
	mock.ExpectRollback()

	req, _ := http.NewRequest("DELETE", "/productos/"+productoID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.Database), response.Error.Code)
	assert.Equal(t, "internal server error", response.Error.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}
