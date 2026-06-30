package service

import (
	"github.com/azmiagr/sakutera-softdev/internal/repository"
	"github.com/azmiagr/sakutera-softdev/pkg/bcrypt"
	"github.com/azmiagr/sakutera-softdev/pkg/jwt"
)

type Service struct {
	AuthService       IAuthService
	OnboardingService IOnboardingService
	JwtAuth           jwt.Interface
}

func NewService(repository *repository.Repository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface) *Service {
	return &Service{
		AuthService:       NewAuthService(repository.UserRepo, repository.SessionRepo, repository.OTPRepo, jwtAuth),
		OnboardingService: NewOnboardingService(repository.UserRepo, repository.WorkCategoryRepo, repository.WorkPlatformRepo),
		JwtAuth:           jwtAuth,
	}
}
