package persistence

import (
	"context"

	"backend-productos/internal/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductoRepositoryImpl struct {
	db *gorm.DB
}

func NewProductoRepository(db *gorm.DB) *ProductoRepositoryImpl {
	return &ProductoRepositoryImpl{db: db}
}

func (r *ProductoRepositoryImpl) GetAll(ctx context.Context, page, limit int) ([]model.Producto, error) {
	var productos []model.Producto
	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&productos).Error
	return productos, err
}

func (r *ProductoRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.Producto, error) {
	var producto model.Producto
	err := r.db.WithContext(ctx).Where("id = ?", id.String()).Take(&producto).Error
	if err != nil {
		return nil, err
	}
	return &producto, nil
}

func (r *ProductoRepositoryImpl) Create(ctx context.Context, p *model.Producto) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *ProductoRepositoryImpl) Update(ctx context.Context, p *model.Producto) error {
	return r.db.WithContext(ctx).Model(p).Updates(p).Error
}

func (r *ProductoRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id.String()).Delete(&model.Producto{}).Error
}
