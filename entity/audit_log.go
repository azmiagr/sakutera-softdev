package entity

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	AuditLogID uuid.UUID  `json:"audit_log_id" gorm:"type:varchar(36);primaryKey"`
	UserID     *uuid.UUID `json:"user_id" gorm:"type:varchar(36);index"`
	DeviceID   *uuid.UUID `json:"device_id" gorm:"type:varchar(36);index"`
	Action     string     `json:"action" gorm:"type:varchar(100);not null"`
	Detail     string     `json:"detail" gorm:"type:varchar(500)"`
	CreatedAt  time.Time  `json:"created_at" gorm:"autoCreateTime"`
}
