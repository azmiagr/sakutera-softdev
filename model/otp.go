package model

import "github.com/google/uuid"

type OTPParam struct {
	OTPID  uuid.UUID `json:"-"`
	UserID uuid.UUID `json:"-"`
	Code   string    `json:"-"`
	Type   string    `json:"-"`
}

type VerifyOTPRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

type VerifyOTPResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}
