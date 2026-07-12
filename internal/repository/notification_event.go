package repository

import (
	"time"

	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type INotificationEventRepository interface {
	ExistsByEventID(tx *gorm.DB, eventID string) (bool, error)
	Create(tx *gorm.DB, e *entity.NotificationEvent) error
	UpdateProcessingStatus(tx *gorm.DB, id uuid.UUID, status, parserVersion string) error
	RedactOldEvents(tx *gorm.DB, olderThan time.Time) (int64, error)
}

type NotificationEventRepository struct {
	db *gorm.DB
}

func NewNotificationEventRepository(db *gorm.DB) INotificationEventRepository {
	return &NotificationEventRepository{db: db}
}

func (r *NotificationEventRepository) ExistsByEventID(tx *gorm.DB, eventID string) (bool, error) {
	var count int64
	err := tx.Model(&entity.NotificationEvent{}).
		Where("event_id = ?", eventID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *NotificationEventRepository) Create(tx *gorm.DB, e *entity.NotificationEvent) error {
	return tx.Create(e).Error
}

func (r *NotificationEventRepository) UpdateProcessingStatus(tx *gorm.DB, id uuid.UUID, status, parserVersion string) error {
	return tx.Model(&entity.NotificationEvent{}).
		Where("notification_event_id = ?", id).
		Updates(map[string]any{
			"processing_status": status,
			"parser_version":    parserVersion,
		}).Error
}

func (r *NotificationEventRepository) RedactOldEvents(tx *gorm.DB, olderThan time.Time) (int64, error) {
	result := tx.Model(&entity.NotificationEvent{}).
		Where("created_at < ? AND (title != '' OR text != '' OR big_text != '')", olderThan).
		Updates(map[string]any{
			"title":    "",
			"text":     "",
			"big_text": "",
		})
	return result.RowsAffected, result.Error
}
