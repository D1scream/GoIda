package models

type AuthRequest struct {
	Login    string `json:"login" validate:"required,min=3"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
