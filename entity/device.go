package entity

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	DeviceID   uuid.UUID  `json:"device_id" gorm:"type:varchar(36);primaryKey"`
	UserID     uuid.UUID  `json:"user_id" gorm:"type:varchar(36);index"`
	DeviceName string     `json:"device_name" gorm:"type:varchar(255)"`
	Platform   string     `json:"platform" gorm:"type:varchar(50)"`
	OSVersion  string     `json:"os_version" gorm:"type:varchar(50)"`
	AppVersion string     `json:"app_version" gorm:"type:varchar(50)"`
	TokenHash  string     `json:"-" gorm:"type:char(64);uniqueIndex;not null"`
	IsActive   bool       `json:"is_active" gorm:"default:true;not null"`
	PairedAt   time.Time  `json:"paired_at" gorm:"autoCreateTime"`
	LastSeenAt time.Time  `json:"last_seen_at"`
	RevokedAt  *time.Time `json:"revoked_at"`
}
