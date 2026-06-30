package entity

import "github.com/google/uuid"

type Notification struct {
	NotificationID uuid.UUID `json:"notification_id" gorm:"type:varchar(36);primaryKey"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:varchar(36);index"`
	Type           string    `json:"type" gorm:"type:enum('info','warning','alert');not null"`
}
