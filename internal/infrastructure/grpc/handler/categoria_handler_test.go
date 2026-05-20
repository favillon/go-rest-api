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

type MockCategoriaRepository struct {
	mock.Mock
}

func (m *MockCategoriaRepository) GetAll(ctx context.Context) ([]model.Categoria, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Categoria), args.Error(1)
}

func (m *MockCategoriaRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Categoria, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Categoria), args.Error(1)
}

func (m *MockCategoriaRepository) Create(ctx context.Context, c *model.Categoria) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockCategoriaRepository) Update(ctx context.Context, c *model.Categoria) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockCategoriaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func categoriaGRPCClient(t *testing.T, repo *MockCategoriaRepository) (pb.CategoriaServiceClient, func()) {
	ctx := context.Background()
	listener := bufconn.Listen(1024 * 1024)

	svc := service.NewCategoriaService(repo)
	s := grpc.NewServer()
	pb.RegisterCategoriaServiceServer(s, handler.NewCategoriaHandler(svc))

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

	return pb.NewCategoriaServiceClient(conn), cleanup
}

func TestCategoriaHandler_Create_Success(t *testing.T) {
	repo := new(MockCategoriaRepository)
	client, cleanup := categoriaGRPCClient(t, repo)
	defer cleanup()

	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Categoria")).Return(nil)

	resp, err := client.CreateCategoria(context.Background(), &pb.CreateCategoriaRequest{
		Nombre:      "Test",
		Descripcion: "Desc",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp.Categoria)
	assert.Equal(t, "Test", resp.Categoria.Nombre)
	repo.AssertExpectations(t)
}

func TestCategoriaHandler_Get_Success(t *testing.T) {
	repo := new(MockCategoriaRepository)
	client, cleanup := categoriaGRPCClient(t, repo)
	defer cleanup()

	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(&model.Categoria{
		ID:          id,
		Nombre:      "Test",
		Descripcion: "Desc",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil)

	resp, err := client.GetCategoria(context.Background(), &pb.GetCategoriaRequest{Id: id.String()})

	assert.NoError(t, err)
	assert.Equal(t, id.String(), resp.Categoria.Id)
	repo.AssertExpectations(t)
}

func TestCategoriaHandler_Get_NotFound(t *testing.T) {
	repo := new(MockCategoriaRepository)
	client, cleanup := categoriaGRPCClient(t, repo)
	defer cleanup()

	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return((*model.Categoria)(nil), domain.ErrNotFound)

	_, err := client.GetCategoria(context.Background(), &pb.GetCategoriaRequest{Id: id.String()})

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	repo.AssertExpectations(t)
}

func TestCategoriaHandler_Update_Success(t *testing.T) {
	repo := new(MockCategoriaRepository)
	client, cleanup := categoriaGRPCClient(t, repo)
	defer cleanup()

	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(&model.Categoria{
		ID:          id,
		Nombre:      "Viejo",
		Descripcion: "Desc",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*model.Categoria")).Return(nil)

	resp, err := client.UpdateCategoria(context.Background(), &pb.UpdateCategoriaRequest{
		Id:          id.String(),
		Nombre:      "Nuevo",
		Descripcion: "Nueva desc",
	})

	assert.NoError(t, err)
	assert.Equal(t, "Nuevo", resp.Categoria.Nombre)
	repo.AssertExpectations(t)
}

func TestCategoriaHandler_Delete_Success(t *testing.T) {
	repo := new(MockCategoriaRepository)
	client, cleanup := categoriaGRPCClient(t, repo)
	defer cleanup()

	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(&model.Categoria{ID: id, Nombre: "Test"}, nil)
	repo.On("Delete", mock.Anything, id).Return(nil)

	_, err := client.DeleteCategoria(context.Background(), &pb.DeleteCategoriaRequest{Id: id.String()})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCategoriaHandler_List_Success(t *testing.T) {
	repo := new(MockCategoriaRepository)
	client, cleanup := categoriaGRPCClient(t, repo)
	defer cleanup()

	repo.On("GetAll", mock.Anything).Return([]model.Categoria{
		{ID: uuid.New(), Nombre: "C1"},
		{ID: uuid.New(), Nombre: "C2"},
	}, nil)

	resp, err := client.ListCategorias(context.Background(), &pb.ListCategoriasRequest{})

	assert.NoError(t, err)
	assert.Len(t, resp.Categorias, 2)
	repo.AssertExpectations(t)
}

func TestCategoriaHandler_List_Error(t *testing.T) {
	repo := new(MockCategoriaRepository)
	client, cleanup := categoriaGRPCClient(t, repo)
	defer cleanup()

	repo.On("GetAll", mock.Anything).Return([]model.Categoria{}, errors.New("db failure"))

	_, err := client.ListCategorias(context.Background(), &pb.ListCategoriasRequest{})

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	repo.AssertExpectations(t)
}
