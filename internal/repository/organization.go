package repository

import (
	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IOrganizationRepository interface {
	GetByID(tx *gorm.DB, orgID uuid.UUID) (*entity.Organization, error)
	GetAll(tx *gorm.DB) ([]entity.Organization, error)
}

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) IOrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) GetByID(tx *gorm.DB, orgID uuid.UUID) (*entity.Organization, error) {
	var org entity.Organization
	err := tx.Where("organization_id = ?", orgID).First(&org).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *OrganizationRepository) GetAll(tx *gorm.DB) ([]entity.Organization, error) {
	var orgs []entity.Organization
	err := tx.Order("name ASC").Find(&orgs).Error
	return orgs, err
}
