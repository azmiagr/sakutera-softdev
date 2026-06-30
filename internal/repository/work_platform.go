package repository

import (
	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IWorkPlatformRepository interface {
	GetByID(tx *gorm.DB, id uuid.UUID) (*entity.WorkPlatform, error)
}

type WorkPlatformRepository struct {
	db *gorm.DB
}

func NewWorkPlatformRepository(db *gorm.DB) IWorkPlatformRepository {
	return &WorkPlatformRepository{db: db}
}

func (r *WorkPlatformRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*entity.WorkPlatform, error) {
	var platform entity.WorkPlatform
	err := tx.Where("work_platform_id = ?", id).First(&platform).Error
	if err != nil {
		return nil, err
	}
	return &platform, nil
}
