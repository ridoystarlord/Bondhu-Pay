package repository

import (
	"context"

	"github.com/ridoystarlord/bondhu-pay/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	*BaseRepository
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository(collection),
	}
}

// FindByPhone finds a user by phone
func (r *UserRepository) FindByPhone(ctx context.Context, mobileNumber string) (*models.User, error) {
	var user models.User
	err := r.Collection.FindOne(ctx, bson.M{"mobileNumber": mobileNumber}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]models.User, error) {
	filter := bson.M{"_id": bson.M{"$in": ids}}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
