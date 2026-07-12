package repository

import (
	"github.com/azmiagr/sakutera-softdev/entity"
	"gorm.io/gorm"
)

type ICollectorConfigRepository interface {
	GetConfig(tx *gorm.DB) (*entity.CollectorConfig, error)
	GetEnabledPackages(tx *gorm.DB) ([]entity.AllowedPackage, error)
	IsPackageAllowed(tx *gorm.DB, packageName string) (bool, error)
}

type CollectorConfigRepository struct {
	db *gorm.DB
}

func NewCollectorConfigRepository(db *gorm.DB) ICollectorConfigRepository {
	return &CollectorConfigRepository{db: db}
}

func (r *CollectorConfigRepository) GetConfig(tx *gorm.DB) (*entity.CollectorConfig, error) {
	var config entity.CollectorConfig
	err := tx.Where("id = ?", 1).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *CollectorConfigRepository) GetEnabledPackages(tx *gorm.DB) ([]entity.AllowedPackage, error) {
	var packages []entity.AllowedPackage
	err := tx.Where("enabled = ?", true).Find(&packages).Error
	if err != nil {
		return nil, err
	}
	return packages, nil
}

func (r *CollectorConfigRepository) IsPackageAllowed(tx *gorm.DB, packageName string) (bool, error) {
	var count int64
	err := tx.Model(&entity.AllowedPackage{}).
		Where("package_name = ? AND enabled = ?", packageName, true).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
