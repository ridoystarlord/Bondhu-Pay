package dto

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone" validate:"required,e164"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Phone    string `json:"phone" validate:"required,e164"`
	Password string `json:"password" validate:"required"`
}
