package repository

import (
	"context"

	"github.com/ridoystarlord/bondhu-pay/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TripMemberPaymentRepository struct {
	base *BaseRepository
}

func NewTripMemberPaymentRepository(coll *mongo.Collection) *TripMemberPaymentRepository {
	return &TripMemberPaymentRepository{
		base: NewBaseRepository(coll),
	}
}

func (r *TripMemberPaymentRepository) Create(ctx context.Context, payment models.TripMemberPayment) (*mongo.InsertOneResult, error) {
	return r.base.Create(ctx, payment)
}

func (r *TripMemberPaymentRepository) FindManyByTrip(ctx context.Context, tripID primitive.ObjectID, results interface{}) error {
	filter := bson.M{"tripId": tripID}
	return r.base.FindMany(ctx, filter, 0, 0, results)
}

func (r *TripMemberPaymentRepository) Update(ctx context.Context, id string, update bson.M) (*mongo.UpdateResult, error) {
	return r.base.Update(ctx, id, update)
}

func (r *TripMemberPaymentRepository) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	return r.base.Delete(ctx, id)
}

func (r *TripMemberPaymentRepository) FindMany(ctx context.Context, filter bson.M, results interface{}) error {
	return r.base.FindMany(ctx, filter, 0, 0, results)
}
