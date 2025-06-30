package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripMember struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    TripID   primitive.ObjectID `bson:"tripId" json:"tripId"`
    UserID   primitive.ObjectID `bson:"userId" json:"userId"`
    Role     string             `bson:"role" json:"role"` // "admin", "member", etc.
    JoinedAt time.Time          `bson:"joinedAt" json:"joinedAt"`
}
