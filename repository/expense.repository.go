package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExpenseRepository struct {
	base *BaseRepository
}

func NewExpenseRepository(coll *mongo.Collection) *ExpenseRepository {
	return &ExpenseRepository{
		base: NewBaseRepository(coll),
	}
}

func (r *ExpenseRepository) Create(ctx context.Context, expense interface{}) (*mongo.InsertOneResult, error) {
	return r.base.Create(ctx, expense)
}

func (r *ExpenseRepository) CreateMany(ctx context.Context, expenses []interface{}) (*mongo.InsertManyResult, error) {
	return r.base.CreateMany(ctx, expenses)
}

func (r *ExpenseRepository) FindByID(ctx context.Context, id string, result interface{}) error {
	return r.base.FindByID(ctx, id, result)
}

func (r *ExpenseRepository) FindByTrip(ctx context.Context, tripID primitive.ObjectID, results interface{}) error {
	filter := bson.M{"tripId": tripID}
	return r.base.FindMany(ctx, filter, 0, 0, results)
}

func (r *ExpenseRepository) Update(ctx context.Context, id string, update bson.M) (*mongo.UpdateResult, error) {
	return r.base.Update(ctx, id, update)
}

func (r *ExpenseRepository) UpdateMany(ctx context.Context, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	return r.base.UpdateMany(ctx, filter, update)
}

func (r *ExpenseRepository) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	return r.base.Delete(ctx, id)
}

func (r *ExpenseRepository) DeleteMany(ctx context.Context, filter bson.M) (*mongo.DeleteResult, error) {
	return r.base.DeleteMany(ctx, filter)
}

func (r *ExpenseRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	return r.base.Count(ctx, filter)
}
