package dto

type ExpenseShareInput struct {
	UserID string  `json:"userId" validate:"required"`
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

type CreateExpenseRequest struct {
	TripID   string              `json:"tripId" validate:"required"`
	Amount   float64             `json:"amount" validate:"required,gt=0"`
	PaidBy   string              `json:"paidBy" validate:"required"`
	Category string              `json:"category" validate:"required,oneof=food transport hotel other"`
	Note     string              `json:"note"`
	Shares   []ExpenseShareInput `json:"shares" validate:"required,dive,required"`
}

type UpdateExpenseRequest struct {
	Amount   float64             `json:"amount,omitempty" validate:"omitempty,gt=0"`
	Category string              `json:"category,omitempty" validate:"omitempty,oneof=food transport hotel other"`
	Note     string              `json:"note,omitempty"`
	Shares   []ExpenseShareInput `json:"shares,omitempty" validate:"omitempty,dive"`
}
