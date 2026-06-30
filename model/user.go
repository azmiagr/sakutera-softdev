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

type CheckPhoneRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type CheckPhoneResponse struct {
	HasPIN       bool   `json:"has_pin"`
	SessionToken string `json:"session_token,omitempty"`
}

type SetPINRequest struct {
	PIN string `json:"pin" binding:"required"`
}

type SetPINResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	PIN         string `json:"pin" binding:"required"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}
