package repository

import (
	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IAccessLogRepository interface {
	Create(tx *gorm.DB, log *entity.AccessLog) error
	GetByPassportID(tx *gorm.DB, passportID uuid.UUID) ([]entity.AccessLog, error)
	GetByPassportIDs(tx *gorm.DB, passportIDs []uuid.UUID) ([]entity.AccessLog, error)
}

type AccessLogRepository struct {
	db *gorm.DB
}

func NewAccessLogRepository(db *gorm.DB) IAccessLogRepository {
	return &AccessLogRepository{db: db}
}

func (r *AccessLogRepository) Create(tx *gorm.DB, log *entity.AccessLog) error {
	return tx.Create(log).Error
}

func (r *AccessLogRepository) GetByPassportID(tx *gorm.DB, passportID uuid.UUID) ([]entity.AccessLog, error) {
	var logs []entity.AccessLog
	err := tx.Where("passport_id = ?", passportID).
		Order("accessed_at DESC").
		Find(&logs).Error
	return logs, err
}

func (r *AccessLogRepository) GetByPassportIDs(tx *gorm.DB, passportIDs []uuid.UUID) ([]entity.AccessLog, error) {
	var logs []entity.AccessLog
	err := tx.Where("passport_id IN ?", passportIDs).
		Order("accessed_at DESC").
		Find(&logs).Error
	return logs, err
}
