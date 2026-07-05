package entity

import (
	"time"

	"github.com/google/uuid"
)

type TokenBlacklist struct {
	BlacklistID uuid.UUID `json:"blacklist_id" gorm:"type:varchar(36);primaryKey"`
	Token       string    `json:"token" gorm:"type:varchar(512);not null;uniqueIndex"`
	ExpiredAt   time.Time `json:"expired_at" gorm:"type:datetime;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}
