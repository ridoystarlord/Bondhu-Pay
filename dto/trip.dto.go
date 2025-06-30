package dto

import "time"

// CreateTripRequest represents payload for creating a trip
type CreateTripRequest struct {
	Name      string    `json:"name" validate:"required,min=3,max=100"`
	StartDate time.Time `json:"startDate" validate:"required"`
	EndDate   time.Time `json:"endDate" validate:"required,gtfield=StartDate"`
	CoverPhoto string   `json:"coverPhoto,omitempty"`
}

// UpdateTripRequest represents payload for updating a trip
type UpdateTripRequest struct {
	Name      string    `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	StartDate time.Time `json:"startDate,omitempty"`
	EndDate   time.Time `json:"endDate,omitempty" validate:"omitempty,gtfield=StartDate"`
	CoverPhoto string   `json:"coverPhoto,omitempty"`
}
