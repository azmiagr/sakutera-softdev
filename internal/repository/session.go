package repository

import (
	"github.com/azmiagr/sakutera-softdev/entity"
	"gorm.io/gorm"
)

type ISessionRepository interface {
	CreateSession(tx *gorm.DB, session *entity.Session) error
	GetSessionByToken(tx *gorm.DB, token string) (*entity.Session, error)
	DeleteSessionByUserID(tx *gorm.DB, userID string) error
}

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) ISessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(tx *gorm.DB, session *entity.Session) error {
	err := tx.Create(session).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *SessionRepository) GetSessionByToken(tx *gorm.DB, token string) (*entity.Session, error) {
	var session entity.Session
	err := tx.Where("token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) DeleteSessionByUserID(tx *gorm.DB, userID string) error {
	err := tx.Where("user_id = ?", userID).Delete(&entity.Session{}).Error
	if err != nil {
		return err
	}
	return nil
}
