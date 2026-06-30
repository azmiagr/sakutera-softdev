package entity

import "github.com/google/uuid"

type OTP struct {
	OTPID  uuid.UUID `json:"otp_id" gorm:"type:varchar(36);primaryKey"`
	UserID uuid.UUID `json:"user_id" gorm:"type:varchar(36);index"`
	Code   string    `json:"code" gorm:"type:varchar(10);not null"`
	Type   string    `json:"type" gorm:"type:enum('registration','reset_pin','verification');not null"`
}
