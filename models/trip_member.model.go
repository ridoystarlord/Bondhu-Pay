package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripMemberRole string

const (
	TripMemberRoleAdmin  TripMemberRole = "admin"
	TripMemberRoleMember TripMemberRole = "member"
)

type TripMember struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	TripID    primitive.ObjectID `bson:"tripId" json:"tripId"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Role      TripMemberRole     `bson:"role" json:"role"` // enum
	JoinedAt  time.Time          `bson:"joinedAt" json:"joinedAt"`
}
