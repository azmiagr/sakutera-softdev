package model

import "github.com/google/uuid"

type UserParam struct {
	UserID      uuid.UUID `json:"-"`
	PhoneNumber string    `json:"-"`
	PINNumber   string    `json:"-"`
}

type RegisterRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	FullName    string `json:"full_name" binding:"required"`
}

type RegisterResponse struct {
	SessionToken string `json:"session_token"`
	PhoneMasked  string `json:"phone_masked"`
	Message      string `json:"message"`
}
