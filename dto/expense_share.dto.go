package dto

type CreateExpenseShareRequest struct {
	ExpenseID string  `json:"expenseId" validate:"required"`
	TripID    string  `json:"tripId" validate:"required"`
	UserID    string  `json:"userId" validate:"required"`
	Amount    float64 `json:"amount" validate:"required,gt=0"`
}

type UpdateExpenseShareRequest struct {
	Settled       *bool   `json:"settled,omitempty"`
	SettledVia    string  `json:"settledVia,omitempty"`
	TripID        string  `json:"tripId,omitempty"`
	TransactionID string  `json:"transactionId,omitempty"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
}
