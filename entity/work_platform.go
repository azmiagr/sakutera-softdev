package entity

import "github.com/google/uuid"

type WorkPlatform struct {
	WorkPlatformID uuid.UUID `json:"work_platform_id" gorm:"type:varchar(36);primaryKey"`
	WorkCategoryID uuid.UUID `json:"work_category_id" gorm:"type:varchar(36)"`
	Name           string    `json:"name" gorm:"type:varchar(100);not null;unique"`

	Users []User `gorm:"foreignKey:WorkPlatformID;references:WorkPlatformID;constraint:onDelete:CASCADE"`
}
