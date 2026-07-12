package entity

import (
	"time"

	"github.com/google/uuid"
)

type RateLimitEvent struct {
	RateLimitEventID uuid.UUID `json:"rate_limit_event_id" gorm:"type:varchar(36);primaryKey"`
	Key              string    `json:"key" gorm:"type:varchar(255);index"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
}
