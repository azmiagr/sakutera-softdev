package entity

import "github.com/google/uuid"

type WorkCategory struct {
	WorkCategoryID uuid.UUID `json:"work_category_id" gorm:"type:varchar(36);primaryKey"`
	Name           string    `json:"name" gorm:"type:varchar(100);not null;unique"`

	WorkPlatforms []WorkPlatform `gorm:"foreignKey:WorkCategoryID;references:WorkCategoryID;constraint:onDelete:CASCADE"`
}
