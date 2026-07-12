package repository

import (
	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IConsentRepository interface {
	Create(tx *gorm.DB, c *entity.Consent) error
	GetByPassportID(tx *gorm.DB, passportID uuid.UUID) ([]entity.Consent, error)
	GetByPassportIDs(tx *gorm.DB, passportIDs []uuid.UUID) ([]entity.Consent, error)
	GetByConsentID(tx *gorm.DB, consentID uuid.UUID) (*entity.Consent, error)
	GetByPassportAndOrg(tx *gorm.DB, passportID, orgID uuid.UUID) (*entity.Consent, error)
	UpdateStatus(tx *gorm.DB, consentID uuid.UUID, status string) error
}

type ConsentRepository struct {
	db *gorm.DB
}

func NewConsentRepository(db *gorm.DB) IConsentRepository {
	return &ConsentRepository{db: db}
}

func (r *ConsentRepository) Create(tx *gorm.DB, c *entity.Consent) error {
	return tx.Create(c).Error
}

func (r *ConsentRepository) GetByPassportID(tx *gorm.DB, passportID uuid.UUID) ([]entity.Consent, error) {
	var consents []entity.Consent
	err := tx.Where("passport_id = ?", passportID).
		Order("granted_at DESC").
		Find(&consents).Error
	return consents, err
}

func (r *ConsentRepository) GetByPassportIDs(tx *gorm.DB, passportIDs []uuid.UUID) ([]entity.Consent, error) {
	var consents []entity.Consent
	err := tx.Where("passport_id IN ?", passportIDs).
		Order("granted_at DESC").
		Find(&consents).Error
	return consents, err
}

func (r *ConsentRepository) GetByConsentID(tx *gorm.DB, consentID uuid.UUID) (*entity.Consent, error) {
	var c entity.Consent
	err := tx.Where("consent_id = ?", consentID).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ConsentRepository) GetByPassportAndOrg(tx *gorm.DB, passportID, orgID uuid.UUID) (*entity.Consent, error) {
	var c entity.Consent
	err := tx.Where("passport_id = ? AND organization_id = ? AND status = 'active'", passportID, orgID).
		First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ConsentRepository) UpdateStatus(tx *gorm.DB, consentID uuid.UUID, status string) error {
	return tx.Model(&entity.Consent{}).
		Where("consent_id = ?", consentID).
		Update("status", status).Error
}
