package repository

import (
	"github.com/azmiagr/sakutera-softdev/entity"
	"gorm.io/gorm"
)

type IAttachmentRepository interface {
	Create(tx *gorm.DB, a *entity.Attachment) error
}

type AttachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) IAttachmentRepository {
	return &AttachmentRepository{db: db}
}

func (r *AttachmentRepository) Create(tx *gorm.DB, a *entity.Attachment) error {
	return tx.Create(a).Error
}
