package handler_test

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"backend-productos/internal/application/service"
	"backend-productos/internal/domain"
	"backend-productos/internal/domain/model"
	"backend-productos/internal/infrastructure/grpc/handler"
	pb "backend-productos/proto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
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

func productoGRPCClient(t *testing.T, repo *MockProductoRepository) (pb.ProductoServiceClient, func()) {
	ctx := context.Background()
	listener := bufconn.Listen(1024 * 1024)

	svc := service.NewProductoService(repo)
	s := grpc.NewServer()
	pb.RegisterProductoServiceServer(s, handler.NewProductoHandler(svc))

	go func() {
		if err := s.Serve(listener); err != nil {
			// expected on stop
		}
	}()

	conn, err := grpc.DialContext(ctx, "passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}

	cleanup := func() {
		conn.Close()
		s.Stop()
	}

	return pb.NewProductoServiceClient(conn), cleanup
}

func TestProductoHandler_Create_Success(t *testing.T) {
	repo := new(MockProductoRepository)
	client, cleanup := productoGRPCClient(t, repo)
	defer cleanup()

	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Producto")).Return(nil)

	resp, err := client.CreateProducto(context.Background(), &pb.CreateProductoRequest{
		Nombre:      "Test",
		Descripcion: "Desc",
		Precio:      10.0,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp.Producto)
	assert.Equal(t, "Test", resp.Producto.Nombre)
	assert.NotEmpty(t, resp.Producto.Id)
	repo.AssertExpectations(t)
}

func TestProductoHandler_Get_Success(t *testing.T) {
	repo := new(MockProductoRepository)
	client, cleanup := productoGRPCClient(t, repo)
	defer cleanup()

	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(&model.Producto{
		ID:          id,
		Nombre:      "Test",
		Descripcion: "Desc",
		Precio:      10.0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil)

	resp, err := client.GetProducto(context.Background(), &pb.GetProductoRequest{Id: id.String()})

	assert.NoError(t, err)
	assert.Equal(t, id.String(), resp.Producto.Id)
	assert.Equal(t, "Test", resp.Producto.Nombre)
	repo.AssertExpectations(t)
}

func TestProductoHandler_Get_NotFound(t *testing.T) {
	repo := new(MockProductoRepository)
	client, cleanup := productoGRPCClient(t, repo)
	defer cleanup()

	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return((*model.Producto)(nil), domain.ErrNotFound)

	_, err := client.GetProducto(context.Background(), &pb.GetProductoRequest{Id: id.String()})

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	repo.AssertExpectations(t)
}

func TestProductoHandler_Update_Success(t *testing.T) {
	repo := new(MockProductoRepository)
	client, cleanup := productoGRPCClient(t, repo)
	defer cleanup()

	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(&model.Producto{
		ID:          id,
		Nombre:      "Viejo",
		Descripcion: "Desc",
		Precio:      10.0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*model.Producto")).Return(nil)

	resp, err := client.UpdateProducto(context.Background(), &pb.UpdateProductoRequest{
		Id:          id.String(),
		Nombre:      "Nuevo",
		Descripcion: "Nueva desc",
		Precio:      20.0,
	})

	assert.NoError(t, err)
	assert.Equal(t, "Nuevo", resp.Producto.Nombre)
	assert.Equal(t, "Nueva desc", resp.Producto.Descripcion)
	assert.Equal(t, 20.0, resp.Producto.Precio)
	repo.AssertExpectations(t)
}

func TestProductoHandler_Update_NotFound(t *testing.T) {
	repo := new(MockProductoRepository)
	client, cleanup := productoGRPCClient(t, repo)
	defer cleanup()

	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return((*model.Producto)(nil), domain.ErrNotFound)

	_, err := client.UpdateProducto(context.Background(), &pb.UpdateProductoRequest{
		Id:     id.String(),
		Nombre: "Nuevo",
	})

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	repo.AssertExpectations(t)
}

func TestProductoHandler_Delete_Success(t *testing.T) {
	repo := new(MockProductoRepository)
	client, cleanup := productoGRPCClient(t, repo)
	defer cleanup()

	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(&model.Producto{ID: id, Nombre: "Test"}, nil)
	repo.On("Delete", mock.Anything, id).Return(nil)

	_, err := client.DeleteProducto(context.Background(), &pb.DeleteProductoRequest{Id: id.String()})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestProductoHandler_Delete_NotFound(t *testing.T) {
	repo := new(MockProductoRepository)
	client, cleanup := productoGRPCClient(t, repo)
	defer cleanup()

	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return((*model.Producto)(nil), domain.ErrNotFound)

	_, err := client.DeleteProducto(context.Background(), &pb.DeleteProductoRequest{Id: id.String()})

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	repo.AssertExpectations(t)
}

func TestProductoHandler_List_Success(t *testing.T) {
	repo := new(MockProductoRepository)
	client, cleanup := productoGRPCClient(t, repo)
	defer cleanup()

	repo.On("GetAll", mock.Anything, 1, 10).Return([]model.Producto{
		{ID: uuid.New(), Nombre: "P1", Precio: 10},
		{ID: uuid.New(), Nombre: "P2", Precio: 20},
	}, nil)

	resp, err := client.ListProductos(context.Background(), &pb.ListProductosRequest{Page: 1, Limit: 10})

	assert.NoError(t, err)
	assert.Len(t, resp.Productos, 2)
	repo.AssertExpectations(t)
}

func TestProductoHandler_List_Error(t *testing.T) {
	repo := new(MockProductoRepository)
	client, cleanup := productoGRPCClient(t, repo)
	defer cleanup()

	repo.On("GetAll", mock.Anything, 1, 10).Return([]model.Producto{}, errors.New("db failure"))

	_, err := client.ListProductos(context.Background(), &pb.ListProductosRequest{Page: 1, Limit: 10})

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	repo.AssertExpectations(t)
}
