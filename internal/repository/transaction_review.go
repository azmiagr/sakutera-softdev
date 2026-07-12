package repository

import (
	"time"

	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITransactionReviewRepository interface {
	Create(tx *gorm.DB, r *entity.TransactionReview) error
	GetByUserID(tx *gorm.DB, userID uuid.UUID, status string, limit int) ([]entity.TransactionReview, error)
	GetByID(tx *gorm.DB, reviewID uuid.UUID) (*entity.TransactionReview, error)
	UpdateStatus(tx *gorm.DB, reviewID uuid.UUID, status string, transactionID *uuid.UUID) error
}

type TransactionReviewRepository struct {
	db *gorm.DB
}

func NewTransactionReviewRepository(db *gorm.DB) ITransactionReviewRepository {
	return &TransactionReviewRepository{db: db}
}

func (r *TransactionReviewRepository) Create(tx *gorm.DB, review *entity.TransactionReview) error {
	return tx.Create(review).Error
}

func (r *TransactionReviewRepository) GetByUserID(tx *gorm.DB, userID uuid.UUID, status string, limit int) ([]entity.TransactionReview, error) {
	var reviews []entity.TransactionReview
	query := tx.Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("created_at DESC").Limit(limit).Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *TransactionReviewRepository) GetByID(tx *gorm.DB, reviewID uuid.UUID) (*entity.TransactionReview, error) {
	var review entity.TransactionReview
	err := tx.Where("review_id = ?", reviewID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *TransactionReviewRepository) UpdateStatus(tx *gorm.DB, reviewID uuid.UUID, status string, transactionID *uuid.UUID) error {
	return tx.Model(&entity.TransactionReview{}).
		Where("review_id = ?", reviewID).
		Updates(map[string]any{
			"status":         status,
			"transaction_id": transactionID,
			"reviewed_at":    time.Now(),
		}).Error
}
