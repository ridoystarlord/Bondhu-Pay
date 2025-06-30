package dto

type CreateTripMemberRequest struct {
	UserID string `json:"userId" validate:"required"`
	Role   string `json:"role" validate:"required,oneof=admin member"`
}

type UpdateTripMemberRequest struct {
	Role string `json:"role" validate:"required,oneof=admin member"`
}
