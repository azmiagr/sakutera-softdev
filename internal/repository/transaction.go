package repository

import (
	"time"

	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITransactionRepository interface {
	CreateTransaction(tx *gorm.DB, t *entity.Transaction) error
	GetLastByUserID(tx *gorm.DB, userID uuid.UUID) (*entity.Transaction, error)
	GetByUserID(tx *gorm.DB, userID uuid.UUID, limit int) ([]entity.Transaction, error)
	GetByUserIDAndDateRange(tx *gorm.DB, userID uuid.UUID, from, to time.Time) ([]entity.Transaction, error)
	CountSuccessByUserID(tx *gorm.DB, userID uuid.UUID) (int64, error)
	SumTodayByUserID(tx *gorm.DB, userID uuid.UUID) (float64, error)
}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) ITransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(tx *gorm.DB, t *entity.Transaction) error {
	return tx.Create(t).Error
}

func (r *TransactionRepository) GetLastByUserID(tx *gorm.DB, userID uuid.UUID) (*entity.Transaction, error) {
	var t entity.Transaction
	err := tx.Where("user_id = ? AND status = ?", userID, "success").
		Order("transaction_date DESC, created_at DESC").
		First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TransactionRepository) GetByUserID(tx *gorm.DB, userID uuid.UUID, limit int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := tx.Where("user_id = ? AND status = ?", userID, "success").
		Order("transaction_date DESC, created_at DESC").
		Limit(limit).
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) GetByUserIDAndDateRange(tx *gorm.DB, userID uuid.UUID, from, to time.Time) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := tx.Where("user_id = ? AND status = ? AND transaction_date BETWEEN ? AND ?",
		userID, "success", from, to).
		Order("transaction_date ASC").
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) CountSuccessByUserID(tx *gorm.DB, userID uuid.UUID) (int64, error) {
	var count int64
	err := tx.Model(&entity.Transaction{}).
		Where("user_id = ? AND status = ?", userID, "success").
		Count(&count).Error
	return count, err
}

func (r *TransactionRepository) SumTodayByUserID(tx *gorm.DB, userID uuid.UUID) (float64, error) {
	today := time.Now().Format("2006-01-02")
	var total float64
	err := tx.Model(&entity.Transaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND status = ? AND DATE(transaction_date) = ?", userID, "success", today).
		Scan(&total).Error
	return total, err
}
