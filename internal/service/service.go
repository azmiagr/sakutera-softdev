package service

import (
	"github.com/azmiagr/sakutera-softdev/internal/repository"
	"github.com/azmiagr/sakutera-softdev/pkg/bcrypt"
	"github.com/azmiagr/sakutera-softdev/pkg/jwt"
	"github.com/azmiagr/sakutera-softdev/pkg/mlclient"
	"github.com/azmiagr/sakutera-softdev/pkg/whatsapp"
)

type Service struct {
	AuthService       IAuthService
	OnboardingService IOnboardingService
	TransactionService ITransactionService
	DashboardService  IDashboardService
	JwtAuth           jwt.Interface
}

func NewService(repository *repository.Repository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface, ml mlclient.Interface, wa whatsapp.Interface) *Service {
	return &Service{
		AuthService:        NewAuthService(repository.UserRepo, repository.SessionRepo, repository.OTPRepo, jwtAuth, wa),
		OnboardingService:  NewOnboardingService(repository.UserRepo, repository.WorkCategoryRepo, repository.WorkPlatformRepo),
		TransactionService: NewTransactionService(repository.TransactionRepo, repository.TransactionSourceRepo, repository.ForecastResultRepo, ml),
		DashboardService:   NewDashboardService(repository.UserRepo, repository.TransactionRepo, repository.TransactionSourceRepo, repository.ForecastResultRepo),
		JwtAuth:            jwtAuth,
	}
}
