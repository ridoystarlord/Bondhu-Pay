package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Expense struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    TripID    primitive.ObjectID `bson:"tripId" json:"tripId"`
    Amount    float64            `bson:"amount" json:"amount"`
    PaidBy    primitive.ObjectID `bson:"paidBy" json:"paidBy"`
    Category  string             `bson:"category" json:"category"` // e.g., "Food", "Transport"
    Note      string             `bson:"note,omitempty" json:"note,omitempty"`
    CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}
