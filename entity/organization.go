package entity

import "github.com/google/uuid"

type Organization struct {
	OrganizationID uuid.UUID `json:"organization_id" gorm:"type:varchar(36);primaryKey"`
	Name           string    `json:"name" gorm:"type:varchar(150);not null"`
	Type           string    `json:"type" gorm:"type:enum('bank','fintech','employer');not null"`

	AccessLogs []AccessLog `gorm:"foreignKey:OrganizationID;references:OrganizationID;constraint:onDelete:CASCADE"`
	Consents   []Consent   `gorm:"foreignKey:OrganizationID;references:OrganizationID;constraint:onDelete:CASCADE"`
}
