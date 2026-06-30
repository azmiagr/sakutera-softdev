package entity

import (
	"time"

	"github.com/google/uuid"
)

type Consent struct {
	ConsentID      uuid.UUID  `json:"consent_id" gorm:"type:varchar(36);primaryKey"`
	PassportID     uuid.UUID  `json:"passport_id" gorm:"type:varchar(36);index"`
	OrganizationID uuid.UUID  `json:"organization_id" gorm:"type:varchar(36);index"`
	Permission     string     `json:"permission" gorm:"type:enum('read','write','full');not null"`
	Status         string     `json:"status" gorm:"type:enum('active','revoked','expired');not null"`
	DataScope      string     `json:"data_scope" gorm:"type:varchar(255);not null;default:'full'"`
	Purpose        string     `json:"purpose" gorm:"type:varchar(255)"`
	GrantedAt      time.Time  `json:"granted_at" gorm:"type:datetime;not null"`
	ExpiresAt      *time.Time `json:"expires_at" gorm:"type:datetime"`
}
