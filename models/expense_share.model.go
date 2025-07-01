package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExpenseShare struct {
    ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
    ExpenseID     primitive.ObjectID `bson:"expenseId" json:"expenseId"`
    UserID        primitive.ObjectID `bson:"userId" json:"userId"`
    Amount        float64            `bson:"amount" json:"amount"`
    Settled       bool               `bson:"settled" json:"settled"`
    SettledVia    string             `bson:"settledVia,omitempty" json:"settledVia,omitempty"` // bKash, Nagad, etc.
    TransactionID string             `bson:"transactionId,omitempty" json:"transactionId,omitempty"`
    SettledAt     *time.Time         `bson:"settledAt,omitempty" json:"settledAt,omitempty"`
}
