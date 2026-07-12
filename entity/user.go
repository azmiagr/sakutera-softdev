package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID           uuid.UUID  `json:"id" gorm:"type:varchar(36);primaryKey"`
	WorkPlatformID   *uuid.UUID `json:"work_platform_id" gorm:"type:varchar(36);index"`
	FullName         string     `json:"full_name" gorm:"type:varchar(255);not null"`
	PhoneNumber      string     `json:"phone_number" gorm:"type:varchar(15);not null"`
	PINNumber        string     `json:"pin_number" gorm:"type:varchar(100)"`
	Status           string     `json:"status" gorm:"type:enum('active', 'inactive');not null;default:'inactive'"`
	IncomeSourceType string     `json:"income_source_type" gorm:"type:varchar(50)"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	Sessions       []Session       `gorm:"foreignKey:UserID;references:UserID;constraint:onDelete:CASCADE"`
	OTPs           []OTP           `gorm:"foreignKey:UserID;references:UserID;constraint:onDelete:CASCADE"`
	Notifications  []Notification  `gorm:"foreignKey:UserID;references:UserID;constraint:onDelete:CASCADE"`
	IncomePassport *IncomePassport `gorm:"foreignKey:UserID;references:UserID;constraint:onDelete:CASCADE"`
	Transactions   []Transaction   `gorm:"foreignKey:UserID;references:UserID;constraint:onDelete:CASCADE"`
	ForecastResult *ForecastResult `gorm:"foreignKey:UserID;references:UserID;constraint:onDelete:CASCADE"`
}
