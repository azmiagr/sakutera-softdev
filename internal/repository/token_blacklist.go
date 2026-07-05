package repository

import (
	"github.com/azmiagr/sakutera-softdev/entity"
	"gorm.io/gorm"
)

type ITokenBlacklistRepository interface {
	CreateBlacklistToken(tx *gorm.DB, blacklist *entity.TokenBlacklist) error
	IsTokenBlacklisted(tx *gorm.DB, token string) (bool, error)
}

type TokenBlacklistRepository struct {
	db *gorm.DB
}

func NewTokenBlacklistRepository(db *gorm.DB) ITokenBlacklistRepository {
	return &TokenBlacklistRepository{db: db}
}

func (r *TokenBlacklistRepository) CreateBlacklistToken(tx *gorm.DB, blacklist *entity.TokenBlacklist) error {
	return tx.Create(blacklist).Error
}

func (r *TokenBlacklistRepository) IsTokenBlacklisted(tx *gorm.DB, token string) (bool, error) {
	var count int64
	err := tx.Model(&entity.TokenBlacklist{}).Where("token = ?", token).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
