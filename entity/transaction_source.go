package entity

import "github.com/google/uuid"

type TransactionSource struct {
	TransactionSourceID uuid.UUID `json:"transaction_source_id" gorm:"type:varchar(36);primaryKey"`
	Name                string    `json:"name" gorm:"type:varchar(255);not null;unique"`
	Provider            string    `json:"provider" gorm:"type:varchar(255)"`
	IsActive            bool      `json:"is_active" gorm:"type:boolean;default:false;not null"`

	Transactions []Transaction `json:"" gorm:"foreignKey:TransactionSourceID;references:TransactionSourceID;constraint:onDelete:CASCADE"`
}
