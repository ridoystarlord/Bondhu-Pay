package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentMethod string

const (
	PaymentMethodCash  PaymentMethod = "cash"
	PaymentMethodBkash PaymentMethod = "bkash"
	PaymentMethodCard  PaymentMethod = "card"
	// Add more as needed
)

type TripMemberPayment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TripID    primitive.ObjectID `bson:"tripId" json:"tripId"`
	MemberID  primitive.ObjectID `bson:"memberId" json:"memberId"`
	Amount    float64            `bson:"amount" json:"amount"`
	Method    PaymentMethod      `bson:"method" json:"method"`
	PaidAt    time.Time          `bson:"paidAt" json:"paidAt"`
	Note      string             `bson:"note,omitempty" json:"note,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
