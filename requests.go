package main

type ValidationRequest struct {
	JWT string `json:"jwt" validate:"required"`
}
