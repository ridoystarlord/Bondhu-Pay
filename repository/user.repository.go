package repository

import (
	"context"

	"github.com/ridoystarlord/bondhu-pay/models"
	"go.mongodb.org/mongo-driver/bson"
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
func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	err := r.Collection.FindOne(ctx, bson.M{"phone": phone}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
