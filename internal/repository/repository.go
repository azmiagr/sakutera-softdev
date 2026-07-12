package repository

import "gorm.io/gorm"

type Repository struct {
	UserRepo              IUserRepository
	SessionRepo           ISessionRepository
	OTPRepo               IOTPRepository
	TokenBlacklistRepo    ITokenBlacklistRepository
	WorkCategoryRepo      IWorkCategoryRepository
	WorkPlatformRepo      IWorkPlatformRepository
	TransactionRepo       ITransactionRepository
	TransactionSourceRepo ITransactionSourceRepository
	ForecastResultRepo    IForecastResultRepository
	IncomePassportRepo    IIncomePassportRepository
	ConsentRepo           IConsentRepository
	AccessLogRepo         IAccessLogRepository
	OrganizationRepo      IOrganizationRepository
	AttachmentRepo        IAttachmentRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepo:              NewUserRepository(db),
		SessionRepo:           NewSessionRepository(db),
		OTPRepo:               NewOTPRepository(db),
		TokenBlacklistRepo:    NewTokenBlacklistRepository(db),
		WorkCategoryRepo:      NewWorkCategoryRepository(db),
		WorkPlatformRepo:      NewWorkPlatformRepository(db),
		TransactionRepo:       NewTransactionRepository(db),
		TransactionSourceRepo: NewTransactionSourceRepository(db),
		ForecastResultRepo:    NewForecastResultRepository(db),
		IncomePassportRepo:    NewIncomePassportRepository(db),
		ConsentRepo:           NewConsentRepository(db),
		AccessLogRepo:         NewAccessLogRepository(db),
		OrganizationRepo:      NewOrganizationRepository(db),
		AttachmentRepo:        NewAttachmentRepository(db),
	}
}
