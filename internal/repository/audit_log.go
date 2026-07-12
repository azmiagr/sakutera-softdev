package repository

import (
	"github.com/azmiagr/sakutera-softdev/entity"
	"gorm.io/gorm"
)

type IAuditLogRepository interface {
	Create(tx *gorm.DB, a *entity.AuditLog) error
}

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) IAuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(tx *gorm.DB, a *entity.AuditLog) error {
	return tx.Create(a).Error
}
