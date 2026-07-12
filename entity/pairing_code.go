package entity

import (
	"time"

	"github.com/google/uuid"
)

type PairingCode struct {
	PairingCodeID uuid.UUID  `json:"pairing_code_id" gorm:"type:varchar(36);primaryKey"`
	UserID        uuid.UUID  `json:"user_id" gorm:"type:varchar(36);index"`
	CodeHash      string     `json:"-" gorm:"type:char(64);uniqueIndex;not null"`
	ExpiresAt     time.Time  `json:"expires_at" gorm:"not null"`
	UsedAt        *time.Time `json:"used_at"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime"`
}
