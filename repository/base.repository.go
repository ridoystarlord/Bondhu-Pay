package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseRepository struct {
	Collection *mongo.Collection
}

func NewBaseRepository(collection *mongo.Collection) *BaseRepository {
	return &BaseRepository{Collection: collection}
}

// FindByID finds one document by ID
func (r *BaseRepository) FindByID(ctx context.Context, id string, result interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}
	filter := bson.M{"_id": objectID}
	return r.Collection.FindOne(ctx, filter).Decode(result)
}

// FindMany finds documents by filter with pagination
func (r *BaseRepository) FindMany(ctx context.Context, filter bson.M, limit int64, skip int64, results interface{}) error {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(limit)
	}
	if skip > 0 {
		opts.SetSkip(skip)
	}

	cursor, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, results)
}

// Create inserts a new document
func (r *BaseRepository) Create(ctx context.Context, doc interface{}) (*mongo.InsertOneResult, error) {
	return r.Collection.InsertOne(ctx, doc)
}

// Update updates a document by ID
func (r *BaseRepository) Update(ctx context.Context, id string, update bson.M) (*mongo.UpdateResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}
	filter := bson.M{"_id": objectID}
	return r.Collection.UpdateOne(ctx, filter, bson.M{"$set": update})
}

// Delete deletes a document by ID
func (r *BaseRepository) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}
	filter := bson.M{"_id": objectID}
	return r.Collection.DeleteOne(ctx, filter)
}

// CreateMany inserts multiple documents at once
func (r *BaseRepository) CreateMany(ctx context.Context, docs []interface{}) (*mongo.InsertManyResult, error) {
	if len(docs) == 0 {
		return nil, errors.New("no documents to insert")
	}
	return r.Collection.InsertMany(ctx, docs)
}

// UpdateMany updates multiple documents matching filter with the update data
func (r *BaseRepository) UpdateMany(ctx context.Context, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	return r.Collection.UpdateMany(ctx, filter, bson.M{"$set": update})
}

// DeleteMany deletes multiple documents matching filter
func (r *BaseRepository) DeleteMany(ctx context.Context, filter bson.M) (*mongo.DeleteResult, error) {
	return r.Collection.DeleteMany(ctx, filter)
}
