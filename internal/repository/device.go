package repository

import (
	"time"

	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IDeviceRepository interface {
	GetByID(tx *gorm.DB, deviceID uuid.UUID) (*entity.Device, error)
	GetByTokenHash(tx *gorm.DB, tokenHash string) (*entity.Device, error)
	GetAllByUserID(tx *gorm.DB, userID uuid.UUID) ([]entity.Device, error)
	Upsert(tx *gorm.DB, d *entity.Device) error
	UpdateLastSeen(tx *gorm.DB, deviceID uuid.UUID) error
	Revoke(tx *gorm.DB, deviceID uuid.UUID) error
}

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) IDeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) GetByID(tx *gorm.DB, deviceID uuid.UUID) (*entity.Device, error) {
	var d entity.Device
	err := tx.Where("device_id = ?", deviceID).First(&d).Error
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DeviceRepository) GetByTokenHash(tx *gorm.DB, tokenHash string) (*entity.Device, error) {
	var d entity.Device
	err := tx.Where("token_hash = ?", tokenHash).First(&d).Error
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DeviceRepository) GetAllByUserID(tx *gorm.DB, userID uuid.UUID) ([]entity.Device, error) {
	var devices []entity.Device
	err := tx.Where("user_id = ?", userID).Order("paired_at DESC").Find(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *DeviceRepository) Upsert(tx *gorm.DB, d *entity.Device) error {
	var existing entity.Device
	err := tx.Where("device_id = ?", d.DeviceID).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return tx.Create(d).Error
	}
	if err != nil {
		return err
	}

	return tx.Model(&entity.Device{}).Where("device_id = ?", d.DeviceID).Updates(map[string]any{
		"user_id":     d.UserID,
		"device_name": d.DeviceName,
		"platform":    d.Platform,
		"os_version":  d.OSVersion,
		"app_version": d.AppVersion,
		"token_hash":  d.TokenHash,
		"is_active":   true,
		"paired_at":   d.PairedAt,
		"revoked_at":  nil,
	}).Error
}

func (r *DeviceRepository) UpdateLastSeen(tx *gorm.DB, deviceID uuid.UUID) error {
	return tx.Model(&entity.Device{}).
		Where("device_id = ?", deviceID).
		Update("last_seen_at", time.Now()).Error
}

func (r *DeviceRepository) Revoke(tx *gorm.DB, deviceID uuid.UUID) error {
	return tx.Model(&entity.Device{}).
		Where("device_id = ?", deviceID).
		Updates(map[string]any{
			"is_active":  false,
			"revoked_at": time.Now(),
		}).Error
}
