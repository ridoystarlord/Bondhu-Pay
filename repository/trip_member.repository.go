package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TripMemberRepository struct {
	base *BaseRepository
}

func NewTripMemberRepository(coll *mongo.Collection) *TripMemberRepository {
	return &TripMemberRepository{
		base: NewBaseRepository(coll),
	}
}

func (r *TripMemberRepository) Create(ctx context.Context, member interface{}) (*mongo.InsertOneResult, error) {
	return r.base.Create(ctx, member)
}

func (r *TripMemberRepository) FindByID(ctx context.Context, id string, result interface{}) error {
	return r.base.FindByID(ctx, id, result)
}

func (r *TripMemberRepository) FindManyByTrip(ctx context.Context, tripID primitive.ObjectID, results interface{}) error {
	filter := bson.M{"tripId": tripID}
	return r.base.FindMany(ctx, filter, 0, 0, results)
}

func (r *TripMemberRepository) Update(ctx context.Context, id string, update bson.M) (*mongo.UpdateResult, error) {
	return r.base.Update(ctx, id, update)
}

func (r *TripMemberRepository) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	return r.base.Delete(ctx, id)
}

func (r *TripMemberRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID, results interface{}) error {
	filter := bson.M{"userId": userID}
	return r.base.FindMany(ctx, filter, 0, 0, results)
}




