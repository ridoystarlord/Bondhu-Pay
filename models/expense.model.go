package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExpenseCategory string

const (
	ExpenseCategoryFood      ExpenseCategory = "food"
	ExpenseCategoryTransport ExpenseCategory = "transport"
	ExpenseCategoryHotel     ExpenseCategory = "hotel"
	ExpenseCategoryOther     ExpenseCategory = "other"
)

type Expense struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TripID    primitive.ObjectID `bson:"tripId" json:"tripId"`
	Amount    float64            `bson:"amount" json:"amount"`
	PaidBy    primitive.ObjectID `bson:"paidBy" json:"paidBy"`
	Category  ExpenseCategory    `bson:"category" json:"category"`
	Note      string             `bson:"note,omitempty" json:"note,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
