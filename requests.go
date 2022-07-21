package main

type AuthRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenRequest struct {
	Username string `json:"username" validate:"required"`
	Id       uint   `json:"id" validate:"required"`
}

type ValidationRequest struct {
	JWT string `json:"jwt" validate:"required"`
}
