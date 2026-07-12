package repository

import (
	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IIncomePassportRepository interface {
	Create(tx *gorm.DB, p *entity.IncomePassport) error
	GetLatestByUserID(tx *gorm.DB, userID uuid.UUID) (*entity.IncomePassport, error)
	GetByID(tx *gorm.DB, passportID uuid.UUID) (*entity.IncomePassport, error)
	GetAllByUserID(tx *gorm.DB, userID uuid.UUID) ([]entity.IncomePassport, error)
}

type IncomePassportRepository struct {
	db *gorm.DB
}

func NewIncomePassportRepository(db *gorm.DB) IIncomePassportRepository {
	return &IncomePassportRepository{db: db}
}

func (r *IncomePassportRepository) Create(tx *gorm.DB, p *entity.IncomePassport) error {
	return tx.Create(p).Error
}

func (r *IncomePassportRepository) GetLatestByUserID(tx *gorm.DB, userID uuid.UUID) (*entity.IncomePassport, error) {
	var p entity.IncomePassport
	err := tx.Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *IncomePassportRepository) GetByID(tx *gorm.DB, passportID uuid.UUID) (*entity.IncomePassport, error) {
	var p entity.IncomePassport
	err := tx.Where("income_passport_id = ?", passportID).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *IncomePassportRepository) GetAllByUserID(tx *gorm.DB, userID uuid.UUID) ([]entity.IncomePassport, error) {
	var passports []entity.IncomePassport
	err := tx.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&passports).Error
	return passports, err
}
