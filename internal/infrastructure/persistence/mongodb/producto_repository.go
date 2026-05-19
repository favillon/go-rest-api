package mongodb

import (
	"context"
	"errors"
	"time"

	"backend-productos/internal/domain"
	"backend-productos/internal/domain/model"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ProductoRepository struct {
	collection *mongo.Collection
}

func NewProductoRepository(db *mongo.Database) *ProductoRepository {
	return &ProductoRepository{
		collection: db.Collection("productos"),
	}
}

func (r *ProductoRepository) GetAll(ctx context.Context, page, limit int) ([]model.Producto, error) {
	skip := (page - 1) * limit

	filter := bson.M{"deleted_at": nil}
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var productos []model.Producto
	if err := cursor.All(ctx, &productos); err != nil {
		return nil, err
	}

	return productos, nil
}

func (r *ProductoRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Producto, error) {
	filter := bson.M{"_id": id, "deleted_at": nil}

	var producto model.Producto
	if err := r.collection.FindOne(ctx, filter).Decode(&producto); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &producto, nil
}

func (r *ProductoRepository) Create(ctx context.Context, p *model.Producto) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, p)
	return err
}

func (r *ProductoRepository) Update(ctx context.Context, p *model.Producto) error {
	p.UpdatedAt = time.Now()

	filter := bson.M{"_id": p.ID, "deleted_at": nil}
	update := bson.M{"$set": p}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *ProductoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	filter := bson.M{"_id": id, "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}
