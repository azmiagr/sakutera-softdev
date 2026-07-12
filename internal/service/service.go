package service

import (
	"github.com/azmiagr/sakutera-softdev/internal/repository"
	"github.com/azmiagr/sakutera-softdev/pkg/bcrypt"
	"github.com/azmiagr/sakutera-softdev/pkg/jwt"
	"github.com/azmiagr/sakutera-softdev/pkg/mlclient"
	"github.com/azmiagr/sakutera-softdev/pkg/supabase"
	"github.com/azmiagr/sakutera-softdev/pkg/whatsapp"
)

type Service struct {
	AuthService              IAuthService
	OnboardingService        IOnboardingService
	TransactionService       ITransactionService
	DashboardService         IDashboardService
	PassportService          IPassportService
	AccessService            IAccessService
	CollectorService         ICollectorService
	TransactionReviewService ITransactionReviewService
	JwtAuth                  jwt.Interface
}

func NewService(repository *repository.Repository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface, ml mlclient.Interface, wa whatsapp.Interface, storage supabase.Interface) *Service {
	return &Service{
		AuthService:              NewAuthService(repository.UserRepo, repository.SessionRepo, repository.OTPRepo, repository.TokenBlacklistRepo, jwtAuth, bcrypt, wa),
		OnboardingService:        NewOnboardingService(repository.UserRepo, repository.WorkCategoryRepo, repository.WorkPlatformRepo),
		TransactionService:       NewTransactionService(repository.TransactionRepo, repository.TransactionSourceRepo, repository.ForecastResultRepo, repository.AttachmentRepo, ml, storage),
		DashboardService:         NewDashboardService(repository.UserRepo, repository.TransactionRepo, repository.TransactionSourceRepo, repository.ForecastResultRepo),
		PassportService:          NewPassportService(repository.IncomePassportRepo, repository.TransactionRepo, repository.ForecastResultRepo),
		AccessService:            NewAccessService(repository.IncomePassportRepo, repository.ConsentRepo, repository.AccessLogRepo, repository.OrganizationRepo),
		CollectorService:         NewCollectorService(repository.DeviceRepo, repository.NotificationEventRepo, repository.PairingCodeRepo, repository.TransactionRepo, repository.TransactionSourceRepo, repository.TransactionReviewRepo, repository.CollectorConfigRepo, repository.AuditLogRepo, repository.RateLimitRepo),
		TransactionReviewService: NewTransactionReviewService(repository.TransactionReviewRepo, repository.TransactionRepo, repository.TransactionSourceRepo, repository.AuditLogRepo),
		JwtAuth:                  jwtAuth,
	}
}
