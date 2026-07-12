package entity

import (
	"time"

	"github.com/google/uuid"
)

type NotificationEvent struct {
	NotificationEventID uuid.UUID `json:"notification_event_id" gorm:"type:varchar(36);primaryKey"`
	EventID             string    `json:"event_id" gorm:"type:char(64);uniqueIndex;not null"`
	UserID              uuid.UUID `json:"user_id" gorm:"type:varchar(36);index"`
	DeviceID            uuid.UUID `json:"device_id" gorm:"type:varchar(36);index"`
	PackageName         string    `json:"package_name" gorm:"type:varchar(255);not null"`
	NotificationID      int64     `json:"notification_id"`
	Title               string    `json:"title" gorm:"type:varchar(500)"`
	Text                string    `json:"text" gorm:"type:varchar(1000)"`
	BigText             string    `json:"big_text" gorm:"type:varchar(2000)"`
	PostedAt            time.Time `json:"posted_at"`
	CapturedAt          time.Time `json:"captured_at"`
	ProcessingStatus    string    `json:"processing_status" gorm:"type:varchar(30);default:'stored';not null"`
	ParserVersion       string    `json:"parser_version" gorm:"type:varchar(30)"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`
}
