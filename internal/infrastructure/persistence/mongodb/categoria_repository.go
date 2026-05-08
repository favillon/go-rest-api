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
)

type CategoriaRepository struct {
	collection *mongo.Collection
}

func NewCategoriaRepository(db *mongo.Database) *CategoriaRepository {
	return &CategoriaRepository{
		collection: db.Collection("categorias"),
	}
}

func (r *CategoriaRepository) GetAll(ctx context.Context) ([]model.Categoria, error) {
	filter := bson.M{"deleted_at": nil}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categorias []model.Categoria
	if err := cursor.All(ctx, &categorias); err != nil {
		return nil, err
	}

	return categorias, nil
}

func (r *CategoriaRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Categoria, error) {
	filter := bson.M{"_id": id.String(), "deleted_at": nil}

	var categoria model.Categoria
	if err := r.collection.FindOne(ctx, filter).Decode(&categoria); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &categoria, nil
}

func (r *CategoriaRepository) Create(ctx context.Context, c *model.Categoria) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, c)
	return err
}

func (r *CategoriaRepository) Update(ctx context.Context, c *model.Categoria) error {
	c.UpdatedAt = time.Now()

	filter := bson.M{"_id": c.ID.String(), "deleted_at": nil}
	update := bson.M{"$set": c}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *CategoriaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	filter := bson.M{"_id": id.String(), "deleted_at": nil}
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
