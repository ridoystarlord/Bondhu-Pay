package dto

type CreateTripMemberPaymentRequest struct {
	MemberID string  `json:"memberId" validate:"required"`          // TripMember._id
	Amount   float64 `json:"amount" validate:"required,gt=0"`       // Must be >0
	Method   string  `json:"method" validate:"required"`            // Enum string: cash/bkash/card
	PaidAt   string  `json:"paidAt,omitempty"`                      // Optional ISO datetime string
	Note     string  `json:"note,omitempty"`
}

type UpdateTripMemberPaymentRequest struct {
	Amount float64 `json:"amount,omitempty" validate:"omitempty,gt=0"`
	Method string  `json:"method,omitempty"`
	PaidAt string  `json:"paidAt,omitempty"`
	Note   string  `json:"note,omitempty"`
}
