package repository

import (
	"time"

	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPairingCodeRepository interface {
	Create(tx *gorm.DB, p *entity.PairingCode) error
	GetActiveByCodeHash(tx *gorm.DB, codeHash string) (*entity.PairingCode, error)
	MarkUsed(tx *gorm.DB, id uuid.UUID) error
	InvalidateActiveByUserID(tx *gorm.DB, userID uuid.UUID) error
	CountRecentByUserID(tx *gorm.DB, userID uuid.UUID, since time.Time) (int64, error)
}

type PairingCodeRepository struct {
	db *gorm.DB
}

func NewPairingCodeRepository(db *gorm.DB) IPairingCodeRepository {
	return &PairingCodeRepository{db: db}
}

func (r *PairingCodeRepository) Create(tx *gorm.DB, p *entity.PairingCode) error {
	return tx.Create(p).Error
}

func (r *PairingCodeRepository) GetActiveByCodeHash(tx *gorm.DB, codeHash string) (*entity.PairingCode, error) {
	var p entity.PairingCode
	err := tx.Where("code_hash = ? AND used_at IS NULL AND expires_at > ?", codeHash, time.Now()).
		First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PairingCodeRepository) MarkUsed(tx *gorm.DB, id uuid.UUID) error {
	return tx.Model(&entity.PairingCode{}).
		Where("pairing_code_id = ?", id).
		Update("used_at", time.Now()).Error
}

func (r *PairingCodeRepository) InvalidateActiveByUserID(tx *gorm.DB, userID uuid.UUID) error {
	return tx.Model(&entity.PairingCode{}).
		Where("user_id = ? AND used_at IS NULL AND expires_at > ?", userID, time.Now()).
		Update("used_at", time.Now()).Error
}

func (r *PairingCodeRepository) CountRecentByUserID(tx *gorm.DB, userID uuid.UUID, since time.Time) (int64, error) {
	var count int64
	err := tx.Model(&entity.PairingCode{}).
		Where("user_id = ? AND created_at > ?", userID, since).
		Count(&count).Error
	return count, err
}
