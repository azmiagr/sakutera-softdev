package repository

import (
	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IWorkCategoryRepository interface {
	GetAllWithPlatforms(tx *gorm.DB) ([]entity.WorkCategory, error)
	GetByID(tx *gorm.DB, id uuid.UUID) (*entity.WorkCategory, error)
}

type WorkCategoryRepository struct {
	db *gorm.DB
}

func NewWorkCategoryRepository(db *gorm.DB) IWorkCategoryRepository {
	return &WorkCategoryRepository{db: db}
}

func (r *WorkCategoryRepository) GetAllWithPlatforms(tx *gorm.DB) ([]entity.WorkCategory, error) {
	var categories []entity.WorkCategory
	err := tx.Preload("WorkPlatforms").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *WorkCategoryRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*entity.WorkCategory, error) {
	var category entity.WorkCategory
	err := tx.Where("work_category_id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}
