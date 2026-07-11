package repository

import (
	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITransactionSourceRepository interface {
	GetAll(tx *gorm.DB, provider string) ([]entity.TransactionSource, error)
	GetByID(tx *gorm.DB, id uuid.UUID) (*entity.TransactionSource, error)
}

type TransactionSourceRepository struct {
	db *gorm.DB
}

func NewTransactionSourceRepository(db *gorm.DB) ITransactionSourceRepository {
	return &TransactionSourceRepository{db: db}
}

func (r *TransactionSourceRepository) GetAll(tx *gorm.DB, provider string) ([]entity.TransactionSource, error) {
	var sources []entity.TransactionSource
	query := tx.Where("is_active = ?", true)
	if provider != "" {
		query = query.Where("provider = ?", provider)
	}
	err := query.Find(&sources).Error
	if err != nil {
		return nil, err
	}
	return sources, nil
}

func (r *TransactionSourceRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*entity.TransactionSource, error) {
	var source entity.TransactionSource
	err := tx.Where("transaction_source_id = ?", id).First(&source).Error
	if err != nil {
		return nil, err
	}
	return &source, nil
}
