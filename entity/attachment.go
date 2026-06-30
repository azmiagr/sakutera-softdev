package entity

import "github.com/google/uuid"

type Attachment struct {
	AttachmentID  uuid.UUID `json:"attachment_id" gorm:"type:varchar(36);primaryKey"`
	TransactionID uuid.UUID `json:"transaction_id" gorm:"type:varchar(36);index"`
	FileURL       string    `json:"file_url" gorm:"type:text;not null"`
	FileType      string    `json:"file_type" gorm:"type:char(50);not null"`
}
