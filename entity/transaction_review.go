package entity

import (
	"time"

	"github.com/google/uuid"
)

type TransactionReview struct {
	ReviewID            uuid.UUID  `json:"review_id" gorm:"type:varchar(36);primaryKey"`
	NotificationEventID uuid.UUID  `json:"notification_event_id" gorm:"type:varchar(36);uniqueIndex"`
	UserID              uuid.UUID  `json:"user_id" gorm:"type:varchar(36);index"`
	Provider            string     `json:"provider" gorm:"type:varchar(100)"`
	TransactionType     string     `json:"transaction_type" gorm:"type:varchar(20)"`
	Amount              float64    `json:"amount" gorm:"type:decimal(10,2)"`
	Description         string     `json:"description" gorm:"type:varchar(500)"`
	TransactionDate     time.Time  `json:"transaction_date"`
	TransactionSourceID *uuid.UUID `json:"transaction_source_id" gorm:"type:varchar(36)"`
	Confidence          float64    `json:"confidence"`
	Reason              string     `json:"reason" gorm:"type:varchar(255)"`
	Status              string     `json:"status" gorm:"type:varchar(20);default:'pending';not null"`
	TransactionID       *uuid.UUID `json:"transaction_id" gorm:"type:varchar(36)"`
	ReviewedAt          *time.Time `json:"reviewed_at"`
	CreatedAt           time.Time  `json:"created_at" gorm:"autoCreateTime"`
}
