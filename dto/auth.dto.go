package dto

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	MobileNumber    string `json:"mobileNumber" validate:"required,e164"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	MobileNumber    string `json:"mobileNumber" validate:"required,e164"`
	Password string `json:"password" validate:"required"`
}
