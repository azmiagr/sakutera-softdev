package repository

import (
	"time"

	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRateLimitRepository interface {
	Create(tx *gorm.DB, key string) error
	CountRecentByKey(tx *gorm.DB, key string, since time.Time) (int64, error)
}

type RateLimitRepository struct {
	db *gorm.DB
}

func NewRateLimitRepository(db *gorm.DB) IRateLimitRepository {
	return &RateLimitRepository{db: db}
}

func (r *RateLimitRepository) Create(tx *gorm.DB, key string) error {
	event := &entity.RateLimitEvent{
		RateLimitEventID: uuid.New(),
		Key:              key,
	}
	return tx.Create(event).Error
}

func (r *RateLimitRepository) CountRecentByKey(tx *gorm.DB, key string, since time.Time) (int64, error) {
	var count int64
	err := tx.Model(&entity.RateLimitEvent{}).
		Where("key = ? AND created_at > ?", key, since).
		Count(&count).Error
	return count, err
}
