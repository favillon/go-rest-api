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

type MockInventarioRepository struct {
	mock.Mock
}

func (m *MockInventarioRepository) GetAll(ctx context.Context) ([]model.Inventario, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Inventario), args.Error(1)
}

func (m *MockInventarioRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Inventario, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Inventario), args.Error(1)
}

func (m *MockInventarioRepository) GetByProductoID(ctx context.Context, productoID string) (*model.Inventario, error) {
	args := m.Called(ctx, productoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Inventario), args.Error(1)
}

func (m *MockInventarioRepository) Create(ctx context.Context, i *model.Inventario) error {
	args := m.Called(ctx, i)
	return args.Error(0)
}

func (m *MockInventarioRepository) Update(ctx context.Context, i *model.Inventario) error {
	args := m.Called(ctx, i)
	return args.Error(0)
}

func (m *MockInventarioRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func inventarioGRPCClient(t *testing.T, repo *MockInventarioRepository) (pb.InventarioServiceClient, func()) {
	ctx := context.Background()
	listener := bufconn.Listen(1024 * 1024)

	svc := service.NewInventarioService(repo)
	s := grpc.NewServer()
	pb.RegisterInventarioServiceServer(s, handler.NewInventarioHandler(svc))

	go func() {
		if err := s.Serve(listener); err != nil {
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

	return pb.NewInventarioServiceClient(conn), cleanup
}

func TestInventarioHandler_Create_Success(t *testing.T) {
	repo := new(MockInventarioRepository)
	client, cleanup := inventarioGRPCClient(t, repo)
	defer cleanup()

	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Inventario")).Return(nil)

	resp, err := client.CreateInventario(context.Background(), &pb.CreateInventarioRequest{
		ProductoId: "pid-1",
		Cantidad:   100,
		Almacen:    "Central",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp.Inventario)
	assert.Equal(t, int32(100), resp.Inventario.Cantidad)
	repo.AssertExpectations(t)
}

func TestInventarioHandler_Get_Success(t *testing.T) {
	repo := new(MockInventarioRepository)
	client, cleanup := inventarioGRPCClient(t, repo)
	defer cleanup()

	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(&model.Inventario{
		ID:         id,
		ProductoID: "pid-1",
		Cantidad:   50,
		Almacen:    "A1",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil)

	resp, err := client.GetInventario(context.Background(), &pb.GetInventarioRequest{Id: id.String()})

	assert.NoError(t, err)
	assert.Equal(t, id.String(), resp.Inventario.Id)
	repo.AssertExpectations(t)
}

func TestInventarioHandler_Get_NotFound(t *testing.T) {
	repo := new(MockInventarioRepository)
	client, cleanup := inventarioGRPCClient(t, repo)
	defer cleanup()

	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return((*model.Inventario)(nil), domain.ErrNotFound)

	_, err := client.GetInventario(context.Background(), &pb.GetInventarioRequest{Id: id.String()})

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	repo.AssertExpectations(t)
}

func TestInventarioHandler_GetByProductoId_Success(t *testing.T) {
	repo := new(MockInventarioRepository)
	client, cleanup := inventarioGRPCClient(t, repo)
	defer cleanup()

	repo.On("GetByProductoID", mock.Anything, "pid-1").Return(&model.Inventario{
		ID:         uuid.New(),
		ProductoID: "pid-1",
		Cantidad:   75,
		Almacen:    "B2",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil)

	resp, err := client.GetInventarioByProductoId(context.Background(), &pb.GetInventarioByProductoIdRequest{ProductoId: "pid-1"})

	assert.NoError(t, err)
	assert.Equal(t, "pid-1", resp.Inventario.ProductoId)
	assert.Equal(t, int32(75), resp.Inventario.Cantidad)
	repo.AssertExpectations(t)
}

func TestInventarioHandler_Update_Success(t *testing.T) {
	repo := new(MockInventarioRepository)
	client, cleanup := inventarioGRPCClient(t, repo)
	defer cleanup()

	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(&model.Inventario{
		ID:         id,
		ProductoID: "pid-1",
		Cantidad:   10,
		Almacen:    "Old",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*model.Inventario")).Return(nil)

	resp, err := client.UpdateInventario(context.Background(), &pb.UpdateInventarioRequest{
		Id:         id.String(),
		ProductoId: "pid-1",
		Cantidad:   99,
		Almacen:    "New",
	})

	assert.NoError(t, err)
	assert.Equal(t, int32(99), resp.Inventario.Cantidad)
	assert.Equal(t, "New", resp.Inventario.Almacen)
	repo.AssertExpectations(t)
}

func TestInventarioHandler_Delete_Success(t *testing.T) {
	repo := new(MockInventarioRepository)
	client, cleanup := inventarioGRPCClient(t, repo)
	defer cleanup()

	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(&model.Inventario{ID: id, ProductoID: "pid-1"}, nil)
	repo.On("Delete", mock.Anything, id).Return(nil)

	_, err := client.DeleteInventario(context.Background(), &pb.DeleteInventarioRequest{Id: id.String()})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestInventarioHandler_List_Success(t *testing.T) {
	repo := new(MockInventarioRepository)
	client, cleanup := inventarioGRPCClient(t, repo)
	defer cleanup()

	repo.On("GetAll", mock.Anything).Return([]model.Inventario{
		{ID: uuid.New(), ProductoID: "pid-1", Cantidad: 10},
		{ID: uuid.New(), ProductoID: "pid-2", Cantidad: 20},
	}, nil)

	resp, err := client.ListInventarios(context.Background(), &pb.ListInventariosRequest{})

	assert.NoError(t, err)
	assert.Len(t, resp.Inventarios, 2)
	repo.AssertExpectations(t)
}

func TestInventarioHandler_List_Error(t *testing.T) {
	repo := new(MockInventarioRepository)
	client, cleanup := inventarioGRPCClient(t, repo)
	defer cleanup()

	repo.On("GetAll", mock.Anything).Return([]model.Inventario{}, errors.New("db failure"))

	_, err := client.ListInventarios(context.Background(), &pb.ListInventariosRequest{})

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	repo.AssertExpectations(t)
}
