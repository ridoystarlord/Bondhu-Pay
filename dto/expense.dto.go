package dto

type CreateExpenseRequest struct {
	TripID   string  `json:"tripId" validate:"required"`
	Amount   float64 `json:"amount" validate:"required,gt=0"`
	PaidBy   string  `json:"paidBy" validate:"required"`
	Category string  `json:"category" validate:"required,oneof=food transport hotel other"`
	Note     string  `json:"note"`
}

type UpdateExpenseRequest struct {
	Amount   float64 `json:"amount,omitempty" validate:"omitempty,gt=0"`
	Category string  `json:"category,omitempty" validate:"omitempty,oneof=food transport hotel other"`
	Note     string  `json:"note,omitempty"`
}
