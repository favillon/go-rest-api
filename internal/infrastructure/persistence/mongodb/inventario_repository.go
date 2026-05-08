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

type InventarioRepository struct {
	collection *mongo.Collection
}

func NewInventarioRepository(db *mongo.Database) *InventarioRepository {
	return &InventarioRepository{
		collection: db.Collection("inventarios"),
	}
}

func (r *InventarioRepository) GetAll(ctx context.Context) ([]model.Inventario, error) {
	filter := bson.M{"deleted_at": nil}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var inventarios []model.Inventario
	if err := cursor.All(ctx, &inventarios); err != nil {
		return nil, err
	}

	return inventarios, nil
}

func (r *InventarioRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Inventario, error) {
	filter := bson.M{"_id": id.String(), "deleted_at": nil}

	var inventario model.Inventario
	if err := r.collection.FindOne(ctx, filter).Decode(&inventario); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &inventario, nil
}

func (r *InventarioRepository) GetByProductoID(ctx context.Context, productoID string) (*model.Inventario, error) {
	filter := bson.M{"producto_id": productoID, "deleted_at": nil}

	var inventario model.Inventario
	if err := r.collection.FindOne(ctx, filter).Decode(&inventario); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &inventario, nil
}

func (r *InventarioRepository) Create(ctx context.Context, i *model.Inventario) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	i.CreatedAt = time.Now()
	i.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, i)
	return err
}

func (r *InventarioRepository) Update(ctx context.Context, i *model.Inventario) error {
	i.UpdatedAt = time.Now()

	filter := bson.M{"_id": i.ID.String(), "deleted_at": nil}
	update := bson.M{"$set": i}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *InventarioRepository) Delete(ctx context.Context, id uuid.UUID) error {
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
