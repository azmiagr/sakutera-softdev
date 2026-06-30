package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	SessionID uuid.UUID `json:"session_id" gorm:"type:varchar(36);primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:varchar(36);index"`
	Token     string    `json:"token" gorm:"type:varchar(512);not null;uniqueIndex"`
	ExpiredAt time.Time `json:"expired_at" gorm:"type:datetime;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
