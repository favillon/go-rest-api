package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apierrors "backend-productos/api/errors"
	"backend-productos/internal/application/service"
	"backend-productos/internal/domain"
	"backend-productos/internal/domain/model"
	httpHandler "backend-productos/internal/infrastructure/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductoRepository struct {
	mock.Mock
}

func (m *MockProductoRepository) GetAll(ctx context.Context, page, limit int) ([]model.Producto, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]model.Producto), args.Error(1)
}

func (m *MockProductoRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Producto, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Producto), args.Error(1)
}

func (m *MockProductoRepository) Create(ctx context.Context, p *model.Producto) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockProductoRepository) Update(ctx context.Context, p *model.Producto) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockProductoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
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

func newTestHandler(mockRepo *MockProductoRepository) *httpHandler.ProductoHandler {
	svc := service.NewProductoService(mockRepo)
	return httpHandler.NewProductoHandler(svc)
}

func TestCrearProducto_Exito(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

	r := gin.Default()
	r.POST("/productos", handler.CrearProducto)

	productoReq := model.Producto{
		Nombre: "Producto Test",
		Precio: 10.50,
	}
	jsonValue, _ := json.Marshal(productoReq)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Producto")).Return(nil)

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
	mockRepo.AssertExpectations(t)
}

func TestCrearProducto_ValidacionFallida(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

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
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

	r := gin.Default()
	r.POST("/productos", handler.CrearProducto)

	productoReq := model.Producto{
		Nombre: "Producto Test",
		Precio: 10.50,
	}
	jsonValue, _ := json.Marshal(productoReq)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Producto")).Return(errors.New("db failure"))

	req, _ := http.NewRequest("POST", "/productos", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.Database), response.Error.Code)
	assert.Equal(t, "internal server error", response.Error.Message)
	mockRepo.AssertExpectations(t)
}

func TestActualizarProducto_Exito(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

	r := gin.Default()
	r.PUT("/productos/:id", handler.ActualizarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440000"
	existente := &model.Producto{ID: uuid.MustParse(productoID), Nombre: "Producto viejo", Descripcion: "Descripcion vieja", Precio: 10.50}

	mockRepo.On("GetByID", mock.Anything, uuid.MustParse(productoID)).Return(existente, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Producto")).Return(nil)

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
	mockRepo.AssertExpectations(t)
}

func TestActualizarProducto_IDInvalido(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

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
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

	r := gin.Default()
	r.PUT("/productos/:id", handler.ActualizarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440001"

	mockRepo.On("GetByID", mock.Anything, uuid.MustParse(productoID)).Return((*model.Producto)(nil), domain.ErrNotFound)

	payload := []byte(`{"nombre":"Producto actualizado","descripcion":"Descripcion nueva","precio":20.75}`)
	req, _ := http.NewRequest("PUT", "/productos/"+productoID, bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.NotFound), response.Error.Code)
	assert.Equal(t, "producto no encontrado", response.Error.Message)
	mockRepo.AssertExpectations(t)
}

func TestActualizarProducto_PayloadInvalido(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

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
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

	r := gin.Default()
	r.PUT("/productos/:id", handler.ActualizarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440003"
	mockRepo.On("GetByID", mock.Anything, uuid.MustParse(productoID)).Return((*model.Producto)(nil), errors.New("db failure"))

	payload := []byte(`{"nombre":"Producto actualizado","descripcion":"Descripcion nueva","precio":20.75}`)
	req, _ := http.NewRequest("PUT", "/productos/"+productoID, bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.Database), response.Error.Code)
	assert.Equal(t, "internal server error", response.Error.Message)
	mockRepo.AssertExpectations(t)
}

func TestActualizarProducto_ErrorDBUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

	r := gin.Default()
	r.PUT("/productos/:id", handler.ActualizarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440004"
	existente := &model.Producto{ID: uuid.MustParse(productoID), Nombre: "Producto viejo", Descripcion: "Descripcion vieja", Precio: 10.50}

	mockRepo.On("GetByID", mock.Anything, uuid.MustParse(productoID)).Return(existente, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Producto")).Return(errors.New("db failure"))

	payload := []byte(`{"nombre":"Producto actualizado","descripcion":"Descripcion nueva","precio":20.75}`)
	req, _ := http.NewRequest("PUT", "/productos/"+productoID, bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.Database), response.Error.Code)
	assert.Equal(t, "internal server error", response.Error.Message)
	mockRepo.AssertExpectations(t)
}

func TestObtenerProductos_Exito(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

	r := gin.Default()
	r.GET("/productos", handler.ObtenerProductos)

	productos := []model.Producto{
		{ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440010"), Nombre: "Teclado", Descripcion: "Mecanico", Precio: 99.99},
		{ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440011"), Nombre: "Mouse", Descripcion: "Inalambrico", Precio: 49.50},
	}

	mockRepo.On("GetAll", mock.Anything, 1, 20).Return(productos, nil)

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
	mockRepo.AssertExpectations(t)
}

func TestObtenerProductos_PaginacionConParametros(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

	r := gin.Default()
	r.GET("/productos", handler.ObtenerProductos)

	productos := []model.Producto{
		{ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440012"), Nombre: "Audifonos", Descripcion: "Noise cancelling", Precio: 120.0},
	}

	mockRepo.On("GetAll", mock.Anything, 2, 1).Return(productos, nil)

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
	mockRepo.AssertExpectations(t)
}

func TestObtenerProductos_PaginacionInvalida(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

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
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

	r := gin.Default()
	r.GET("/productos/:id", handler.ObtenerProductoPorID)

	productoID := "550e8400-e29b-41d4-a716-446655440020"
	producto := &model.Producto{ID: uuid.MustParse(productoID), Nombre: "Monitor", Descripcion: "4K", Precio: 299.90}

	mockRepo.On("GetByID", mock.Anything, uuid.MustParse(productoID)).Return(producto, nil)

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
	mockRepo.AssertExpectations(t)
}

func TestObtenerProductoPorID_IDInvalido(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

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
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

	r := gin.Default()
	r.GET("/productos/:id", handler.ObtenerProductoPorID)

	productoID := "550e8400-e29b-41d4-a716-446655440021"

	mockRepo.On("GetByID", mock.Anything, uuid.MustParse(productoID)).Return((*model.Producto)(nil), domain.ErrNotFound)

	req, _ := http.NewRequest("GET", "/productos/"+productoID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.NotFound), response.Error.Code)
	assert.Equal(t, "producto no encontrado", response.Error.Message)
	mockRepo.AssertExpectations(t)
}

func TestObtenerProductos_ErrorDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

	r := gin.Default()
	r.GET("/productos", handler.ObtenerProductos)

	mockRepo.On("GetAll", mock.Anything, 1, 20).Return([]model.Producto{}, errors.New("db failure"))

	req, _ := http.NewRequest("GET", "/productos", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.Database), response.Error.Code)
	assert.Equal(t, "internal server error", response.Error.Message)
	mockRepo.AssertExpectations(t)
}

func TestObtenerProductoPorID_ErrorDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

	r := gin.Default()
	r.GET("/productos/:id", handler.ObtenerProductoPorID)

	productoID := "550e8400-e29b-41d4-a716-446655440022"
	mockRepo.On("GetByID", mock.Anything, uuid.MustParse(productoID)).Return((*model.Producto)(nil), errors.New("db failure"))

	req, _ := http.NewRequest("GET", "/productos/"+productoID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.Database), response.Error.Code)
	assert.Equal(t, "internal server error", response.Error.Message)
	mockRepo.AssertExpectations(t)
}

func TestEliminarProducto_Exito(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

	r := gin.Default()
	r.DELETE("/productos/:id", handler.EliminarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440030"
	existente := &model.Producto{ID: uuid.MustParse(productoID), Nombre: "Camara", Descripcion: "HD", Precio: 149.99}

	mockRepo.On("GetByID", mock.Anything, uuid.MustParse(productoID)).Return(existente, nil)
	mockRepo.On("Delete", mock.Anything, uuid.MustParse(productoID)).Return(nil)

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
	mockRepo.AssertExpectations(t)
}

func TestEliminarProducto_IDInvalido(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

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
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

	r := gin.Default()
	r.DELETE("/productos/:id", handler.EliminarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440031"
	mockRepo.On("GetByID", mock.Anything, uuid.MustParse(productoID)).Return((*model.Producto)(nil), domain.ErrNotFound)

	req, _ := http.NewRequest("DELETE", "/productos/"+productoID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.NotFound), response.Error.Code)
	assert.Equal(t, "producto no encontrado", response.Error.Message)
	mockRepo.AssertExpectations(t)
}

func TestEliminarProducto_ErrorDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

	r := gin.Default()
	r.DELETE("/productos/:id", handler.EliminarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440032"
	mockRepo.On("GetByID", mock.Anything, uuid.MustParse(productoID)).Return((*model.Producto)(nil), errors.New("db failure"))

	req, _ := http.NewRequest("DELETE", "/productos/"+productoID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.Database), response.Error.Code)
	assert.Equal(t, "internal server error", response.Error.Message)
	mockRepo.AssertExpectations(t)
}

func TestEliminarProducto_ErrorDBDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockProductoRepository)
	handler := newTestHandler(mockRepo)

	r := gin.Default()
	r.DELETE("/productos/:id", handler.EliminarProducto)

	productoID := "550e8400-e29b-41d4-a716-446655440033"
	existente := &model.Producto{ID: uuid.MustParse(productoID), Nombre: "Camara", Descripcion: "HD", Precio: 149.99}

	mockRepo.On("GetByID", mock.Anything, uuid.MustParse(productoID)).Return(existente, nil)
	mockRepo.On("Delete", mock.Anything, uuid.MustParse(productoID)).Return(errors.New("db failure"))

	req, _ := http.NewRequest("DELETE", "/productos/"+productoID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := decodeAPIResponse(t, w.Body.Bytes())

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, string(apierrors.Database), response.Error.Code)
	assert.Equal(t, "internal server error", response.Error.Message)
	mockRepo.AssertExpectations(t)
}
