package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Trip struct {
    ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name       string             `bson:"name" json:"name"`
    StartDate  time.Time          `bson:"startDate" json:"startDate"`
    EndDate    time.Time          `bson:"endDate" json:"endDate"`
    CoverPhoto string             `bson:"coverPhoto,omitempty" json:"coverPhoto,omitempty"`
    CreatedByID  primitive.ObjectID `bson:"createdById" json:"createdById"`
    CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
    UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}
