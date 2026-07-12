package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	TransactionID             uuid.UUID  `json:"transaction_id" gorm:"type:varchar(36);primaryKey"`
	UserID                    uuid.UUID  `json:"user_id" gorm:"type:varchar(36);index"`
	TransactionSourceID       uuid.UUID  `json:"transaction_source_id" gorm:"type:varchar(36);index"`
	Amount                    float64    `json:"amount" gorm:"type:decimal(10,2);not null"`
	TransactionDate           time.Time  `json:"transaction_date" gorm:"type:datetime;not null"`
	Category                  string     `json:"category" gorm:"type:varchar(255);not null"`
	ExternalTransactionID     string     `json:"external_transaction_id" gorm:"type:varchar(150)"`
	Status                    string     `json:"status" gorm:"type:enum('pending','success','failed');not null;default:'pending'"`
	PreviousHash              string     `json:"previous_hash" gorm:"type:char(64);not null"`
	CurrentHash               string     `json:"current_hash" gorm:"type:char(64);not null"`
	SourceNotificationEventID *uuid.UUID `json:"source_notification_event_id" gorm:"type:varchar(36)"`
	CreatedAt                 time.Time  `json:"created_at" gorm:"autoCreateTime"`

	Attachments []Attachment `gorm:"foreignKey:TransactionID;references:TransactionID;constraint:onDelete:CASCADE"`
}
