package repository

import (
	"context"

	"github.com/ridoystarlord/bondhu-pay/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TripRepository struct {
	*BaseRepository
}

func NewTripRepository(collection *mongo.Collection) *TripRepository {
	return &TripRepository{
		BaseRepository: NewBaseRepository(collection),
	}
}

// FindByCreatedBy retrieves trips created by a specific user
func (r *TripRepository) FindByCreatedBy(ctx context.Context, userID string) ([]models.Trip, error) {
	filter := bson.M{"createdBy": userID}
	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var trips []models.Trip
	if err := cursor.All(ctx, &trips); err != nil {
		return nil, err
	}
	return trips, nil
}
