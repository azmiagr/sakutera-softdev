package entity

import "time"

type CollectorConfig struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Mode          string    `json:"mode" gorm:"type:varchar(30);not null;default:'whitelist_only'"`
	ConfigVersion int       `json:"config_version" gorm:"not null;default:1"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AllowedPackage struct {
	PackageName string `json:"package_name" gorm:"type:varchar(255);primaryKey"`
	Provider    string `json:"provider" gorm:"type:varchar(100)"`
	Enabled     bool   `json:"enabled" gorm:"default:true;not null"`
}
