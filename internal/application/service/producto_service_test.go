package service_test

import (
	"context"
	"errors"
	"testing"

	"backend-productos/internal/application/service"
	"backend-productos/internal/domain/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
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

func TestService_GetAll_Success(t *testing.T) {
	mockRepo := new(MockProductoRepository)
	svc := service.NewProductoService(mockRepo)

	ctx := context.Background()
	productos := []model.Producto{
		{ID: uuid.New(), Nombre: "Producto 1", Precio: 10.0},
		{ID: uuid.New(), Nombre: "Producto 2", Precio: 20.0},
	}

	mockRepo.On("GetAll", ctx, 1, 10).Return(productos, nil)

	result, err := svc.GetAll(ctx, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Producto 1", result[0].Nombre)
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_Error(t *testing.T) {
	mockRepo := new(MockProductoRepository)
	svc := service.NewProductoService(mockRepo)

	ctx := context.Background()
	mockRepo.On("GetAll", ctx, 1, 10).Return([]model.Producto{}, errors.New("db failure"))

	result, err := svc.GetAll(ctx, 1, 10)

	assert.Error(t, err)
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}

func TestService_GetByID_Success(t *testing.T) {
	mockRepo := new(MockProductoRepository)
	svc := service.NewProductoService(mockRepo)

	ctx := context.Background()
	id := uuid.New()
	producto := &model.Producto{ID: id, Nombre: "Producto Test", Precio: 15.0}

	mockRepo.On("GetByID", ctx, id).Return(producto, nil)

	result, err := svc.GetByID(ctx, id)

	assert.NoError(t, err)
	assert.Equal(t, "Producto Test", result.Nombre)
	mockRepo.AssertExpectations(t)
}

func TestService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockProductoRepository)
	svc := service.NewProductoService(mockRepo)

	ctx := context.Background()
	id := uuid.New()

	mockRepo.On("GetByID", ctx, id).Return((*model.Producto)(nil), gorm.ErrRecordNotFound)

	result, err := svc.GetByID(ctx, id)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestService_GetByID_Error(t *testing.T) {
	mockRepo := new(MockProductoRepository)
	svc := service.NewProductoService(mockRepo)

	ctx := context.Background()
	id := uuid.New()

	mockRepo.On("GetByID", ctx, id).Return((*model.Producto)(nil), errors.New("db failure"))

	result, err := svc.GetByID(ctx, id)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestService_Create_Success(t *testing.T) {
	mockRepo := new(MockProductoRepository)
	svc := service.NewProductoService(mockRepo)

	ctx := context.Background()
	producto := &model.Producto{Nombre: "Nuevo", Precio: 5.0}

	mockRepo.On("Create", ctx, producto).Return(nil)

	err := svc.Create(ctx, producto)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_Create_Error(t *testing.T) {
	mockRepo := new(MockProductoRepository)
	svc := service.NewProductoService(mockRepo)

	ctx := context.Background()
	producto := &model.Producto{Nombre: "Nuevo", Precio: 5.0}

	mockRepo.On("Create", ctx, producto).Return(errors.New("db failure"))

	err := svc.Create(ctx, producto)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_Update_Success(t *testing.T) {
	mockRepo := new(MockProductoRepository)
	svc := service.NewProductoService(mockRepo)

	ctx := context.Background()
	id := uuid.New()
	existente := &model.Producto{ID: id, Nombre: "Viejo", Descripcion: "Desc vieja", Precio: 10.0}
	input := &model.Producto{Nombre: "Nuevo", Descripcion: "Desc nueva", Precio: 20.0}

	mockRepo.On("GetByID", ctx, id).Return(existente, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*model.Producto")).Return(nil)

	result, err := svc.Update(ctx, id, input)

	assert.NoError(t, err)
	assert.Equal(t, "Nuevo", result.Nombre)
	assert.Equal(t, "Desc nueva", result.Descripcion)
	assert.Equal(t, 20.0, result.Precio)
	mockRepo.AssertExpectations(t)
}

func TestService_Update_NotFound(t *testing.T) {
	mockRepo := new(MockProductoRepository)
	svc := service.NewProductoService(mockRepo)

	ctx := context.Background()
	id := uuid.New()
	input := &model.Producto{Nombre: "Nuevo", Precio: 20.0}

	mockRepo.On("GetByID", ctx, id).Return((*model.Producto)(nil), gorm.ErrRecordNotFound)

	result, err := svc.Update(ctx, id, input)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestService_Update_Error(t *testing.T) {
	mockRepo := new(MockProductoRepository)
	svc := service.NewProductoService(mockRepo)

	ctx := context.Background()
	id := uuid.New()
	existente := &model.Producto{ID: id, Nombre: "Viejo", Descripcion: "Desc", Precio: 10.0}
	input := &model.Producto{Nombre: "Nuevo", Precio: 20.0}

	mockRepo.On("GetByID", ctx, id).Return(existente, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*model.Producto")).Return(errors.New("db failure"))

	result, err := svc.Update(ctx, id, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestService_Delete_Success(t *testing.T) {
	mockRepo := new(MockProductoRepository)
	svc := service.NewProductoService(mockRepo)

	ctx := context.Background()
	id := uuid.New()
	existente := &model.Producto{ID: id, Nombre: "Test", Precio: 10.0}

	mockRepo.On("GetByID", ctx, id).Return(existente, nil)
	mockRepo.On("Delete", ctx, id).Return(nil)

	err := svc.Delete(ctx, id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_Delete_NotFound(t *testing.T) {
	mockRepo := new(MockProductoRepository)
	svc := service.NewProductoService(mockRepo)

	ctx := context.Background()
	id := uuid.New()

	mockRepo.On("GetByID", ctx, id).Return((*model.Producto)(nil), gorm.ErrRecordNotFound)

	err := svc.Delete(ctx, id)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	mockRepo.AssertExpectations(t)
}

func TestService_Delete_Error(t *testing.T) {
	mockRepo := new(MockProductoRepository)
	svc := service.NewProductoService(mockRepo)

	ctx := context.Background()
	id := uuid.New()
	existente := &model.Producto{ID: id, Nombre: "Test", Precio: 10.0}

	mockRepo.On("GetByID", ctx, id).Return(existente, nil)
	mockRepo.On("Delete", ctx, id).Return(errors.New("db failure"))

	err := svc.Delete(ctx, id)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
