package persistence_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"backend-productos/internal/domain/model"
	"backend-productos/internal/infrastructure/persistence"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

func TestRepository_GetAll_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := persistence.NewProductoRepository(db)

	now := time.Now()
	mock.ExpectQuery(`SELECT .* FROM "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "descripcion", "precio", "created_at", "updated_at", "deleted_at"}).
			AddRow(uuid.New().String(), "Teclado", "Mecanico", 99.99, now, now, nil))

	result, err := repo.GetAll(context.Background(), 1, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Teclado", result[0].Nombre)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetAll_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := persistence.NewProductoRepository(db)

	mock.ExpectQuery(`SELECT .* FROM "productos"`).WillReturnError(errors.New("db failure"))

	result, err := repo.GetAll(context.Background(), 1, 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := persistence.NewProductoRepository(db)

	id := uuid.New()
	now := time.Now()
	mock.ExpectQuery(`SELECT .* FROM "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "descripcion", "precio", "created_at", "updated_at", "deleted_at"}).
			AddRow(id.String(), "Monitor", "4K", 299.90, now, now, nil))

	result, err := repo.GetByID(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, "Monitor", result.Nombre)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := persistence.NewProductoRepository(db)

	id := uuid.New()
	mock.ExpectQuery(`SELECT .* FROM "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "descripcion", "precio", "created_at", "updated_at", "deleted_at"}))

	result, err := repo.GetByID(context.Background(), id)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := persistence.NewProductoRepository(db)

	id := uuid.New()
	mock.ExpectQuery(`SELECT .* FROM "productos"`).WillReturnError(errors.New("db failure"))

	result, err := repo.GetByID(context.Background(), id)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Create_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := persistence.NewProductoRepository(db)

	producto := &model.Producto{Nombre: "Nuevo", Precio: 10.0}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "productos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New().String()))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), producto)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Create_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := persistence.NewProductoRepository(db)

	producto := &model.Producto{Nombre: "Nuevo", Precio: 10.0}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "productos"`).WillReturnError(errors.New("db failure"))
	mock.ExpectRollback()

	err := repo.Create(context.Background(), producto)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := persistence.NewProductoRepository(db)

	producto := &model.Producto{ID: uuid.New(), Nombre: "Actualizado", Precio: 20.0}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "productos"`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), producto)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := persistence.NewProductoRepository(db)

	producto := &model.Producto{ID: uuid.New(), Nombre: "Actualizado", Precio: 20.0}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "productos"`).WillReturnError(errors.New("db failure"))
	mock.ExpectRollback()

	err := repo.Update(context.Background(), producto)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := persistence.NewProductoRepository(db)

	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "productos"`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), id)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := persistence.NewProductoRepository(db)

	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "productos"`).WillReturnError(errors.New("db failure"))
	mock.ExpectRollback()

	err := repo.Delete(context.Background(), id)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
