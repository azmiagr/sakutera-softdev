package entity

import "github.com/google/uuid"

type Session struct {
	SessionID uuid.UUID `json:"session_id" gorm:"type:varchar(36);primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:varchar(36);index"`
	Token     string    `json:"token" gorm:"type:varchar(255);not null"`
}
