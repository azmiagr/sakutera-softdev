package repository

import (
	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/azmiagr/sakutera-softdev/model"
	"gorm.io/gorm"
)

type IOTPRepository interface {
	CreateOTP(tx *gorm.DB, otp *entity.OTP) error
	GetOTP(tx *gorm.DB, param model.OTPParam) (*entity.OTP, error)
	DeleteOTPByUserID(tx *gorm.DB, userID string) error
}

type OTPRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) IOTPRepository {
	return &OTPRepository{db: db}
}

func (r *OTPRepository) CreateOTP(tx *gorm.DB, otp *entity.OTP) error {
	err := tx.Create(otp).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *OTPRepository) GetOTP(tx *gorm.DB, param model.OTPParam) (*entity.OTP, error) {
	var otp entity.OTP
	err := tx.Where(&param).First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *OTPRepository) DeleteOTPByUserID(tx *gorm.DB, userID string) error {
	err := tx.Where("user_id = ?", userID).Delete(&entity.OTP{}).Error
	if err != nil {
		return err
	}
	return nil
}
