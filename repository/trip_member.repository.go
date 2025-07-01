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

func (r *TripMemberRepository) FindMembersWithUserInfo(ctx context.Context, tripID primitive.ObjectID) ([]bson.M, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "tripId", Value: tripID}}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "userId"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "user"},
		}}},
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$user"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "user.passwordHash", Value: 0},
		}}},
	}

	cursor, err := r.base.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
