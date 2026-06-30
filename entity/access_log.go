package entity

import (
	"time"

	"github.com/google/uuid"
)

type AccessLog struct {
	AccessLogID    uuid.UUID `json:"access_log_id" gorm:"type:varchar(36);primaryKey"`
	PassportID     uuid.UUID `json:"passport_id" gorm:"type:varchar(36);index"`
	OrganizationID uuid.UUID `json:"organization_id" gorm:"type:varchar(36);index"`
	AccessedAt     time.Time `json:"accessed_at" gorm:"type:datetime;not null"`
	Status         string    `json:"status" gorm:"type:enum('success','failed','pending');not null"`
}
