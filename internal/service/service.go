package service

import (
	"github.com/azmiagr/sakutera-softdev/internal/repository"
	"github.com/azmiagr/sakutera-softdev/pkg/bcrypt"
	"github.com/azmiagr/sakutera-softdev/pkg/jwt"
	"github.com/azmiagr/sakutera-softdev/pkg/mlclient"
	"github.com/azmiagr/sakutera-softdev/pkg/whatsapp"
)

type Service struct {
	AuthService        IAuthService
	OnboardingService  IOnboardingService
	TransactionService ITransactionService
	DashboardService   IDashboardService
	PassportService    IPassportService
	AccessService      IAccessService
	JwtAuth            jwt.Interface
}

func NewService(repository *repository.Repository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface, ml mlclient.Interface, wa whatsapp.Interface) *Service {
	return &Service{
		AuthService:        NewAuthService(repository.UserRepo, repository.SessionRepo, repository.OTPRepo, repository.TokenBlacklistRepo, jwtAuth, bcrypt, wa),
		OnboardingService:  NewOnboardingService(repository.UserRepo, repository.WorkCategoryRepo, repository.WorkPlatformRepo),
		TransactionService: NewTransactionService(repository.TransactionRepo, repository.TransactionSourceRepo, repository.ForecastResultRepo, ml),
		DashboardService:   NewDashboardService(repository.UserRepo, repository.TransactionRepo, repository.TransactionSourceRepo, repository.ForecastResultRepo),
		PassportService:    NewPassportService(repository.IncomePassportRepo, repository.TransactionRepo, repository.ForecastResultRepo),
		AccessService:      NewAccessService(repository.IncomePassportRepo, repository.ConsentRepo, repository.AccessLogRepo, repository.OrganizationRepo),
		JwtAuth:            jwtAuth,
	}
}
