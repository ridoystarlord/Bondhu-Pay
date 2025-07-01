package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExpenseShareRepository struct {
	base *BaseRepository
}

func NewExpenseShareRepository(coll *mongo.Collection) *ExpenseShareRepository {
	return &ExpenseShareRepository{
		base: NewBaseRepository(coll),
	}
}

func (r *ExpenseShareRepository) Create(ctx context.Context, share interface{}) (*mongo.InsertOneResult, error) {
	return r.base.Create(ctx, share)
}

func (r *ExpenseShareRepository) CreateMany(ctx context.Context, shares []interface{}) (*mongo.InsertManyResult, error) {
	return r.base.CreateMany(ctx, shares)
}

func (r *ExpenseShareRepository) FindByExpense(ctx context.Context, expenseID primitive.ObjectID, results interface{}) error {
	filter := bson.M{"expenseId": expenseID}
	return r.base.FindMany(ctx, filter, 0, 0, results)
}

func (r *ExpenseShareRepository) Update(ctx context.Context, id string, update bson.M) (*mongo.UpdateResult, error) {
	return r.base.Update(ctx, id, update)
}

func (r *ExpenseShareRepository) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	return r.base.Delete(ctx, id)
}

func (r *ExpenseShareRepository) FindMany(ctx context.Context, filter bson.M, results interface{}) error {
	return r.base.FindMany(ctx, filter, 0, 0, results)
}

func (r *ExpenseShareRepository) FindManyByTrip(ctx context.Context, tripID primitive.ObjectID, results interface{}) error {
	filter := bson.M{"tripId": tripID}
	return r.base.FindMany(ctx, filter, 0, 0, results)
}
