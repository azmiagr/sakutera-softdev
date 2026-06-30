package repository

import "gorm.io/gorm"

type Repository struct {
	UserRepo         IUserRepository
	SessionRepo      ISessionRepository
	OTPRepo          IOTPRepository
	WorkCategoryRepo IWorkCategoryRepository
	WorkPlatformRepo IWorkPlatformRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepo:         NewUserRepository(db),
		SessionRepo:      NewSessionRepository(db),
		OTPRepo:          NewOTPRepository(db),
		WorkCategoryRepo: NewWorkCategoryRepository(db),
		WorkPlatformRepo: NewWorkPlatformRepository(db),
	}
}
