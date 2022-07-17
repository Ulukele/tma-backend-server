package main

type AuthRequest struct {
	Username string `json:"username" validate:"required"`
	Id       uint   `json:"id"`
	Password string `json:"password" validate:"required"`
}

type ValidationRequest struct {
	JWT string `json:"jwt" validate:"required"`
}
