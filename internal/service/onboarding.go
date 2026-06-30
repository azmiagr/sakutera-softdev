package service

import (
	"github.com/azmiagr/sakutera-softdev/internal/repository"
	"github.com/azmiagr/sakutera-softdev/model"
	apperrors "github.com/azmiagr/sakutera-softdev/pkg/errors"
	"github.com/azmiagr/sakutera-softdev/pkg/database/mariadb"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IOnboardingService interface {
	GetWorkCategories() (*model.GetWorkCategoriesResponse, error)
	SelectWorkPlatform(userID uuid.UUID, req model.SelectWorkPlatformRequest) (*model.SelectWorkPlatformResponse, error)
	GetIncomeSources() (*model.GetIncomeSourcesResponse, error)
	SelectIncomeSource(userID uuid.UUID, req model.SelectIncomeSourceRequest) (*model.SelectIncomeSourceResponse, error)
}

type OnboardingService struct {
	db               *gorm.DB
	userRepo         repository.IUserRepository
	workCategoryRepo repository.IWorkCategoryRepository
	workPlatformRepo repository.IWorkPlatformRepository
}

func NewOnboardingService(
	userRepo repository.IUserRepository,
	workCategoryRepo repository.IWorkCategoryRepository,
	workPlatformRepo repository.IWorkPlatformRepository,
) IOnboardingService {
	return &OnboardingService{
		db:               mariadb.Connection,
		userRepo:         userRepo,
		workCategoryRepo: workCategoryRepo,
		workPlatformRepo: workPlatformRepo,
	}
}

func (s *OnboardingService) GetWorkCategories() (*model.GetWorkCategoriesResponse, error) {
	categories, err := s.workCategoryRepo.GetAllWithPlatforms(s.db)
	if err != nil {
		return nil, apperrors.InternalServer("gagal mengambil data pekerjaan")
	}

	var items []model.WorkCategoryItem
	for _, cat := range categories {
		var platforms []model.WorkPlatformItem
		for _, p := range cat.WorkPlatforms {
			platforms = append(platforms, model.WorkPlatformItem{
				WorkPlatformID: p.WorkPlatformID.String(),
				Name:           p.Name,
			})
		}
		items = append(items, model.WorkCategoryItem{
			WorkCategoryID: cat.WorkCategoryID.String(),
			Name:           cat.Name,
			Platforms:      platforms,
		})
	}

	return &model.GetWorkCategoriesResponse{Categories: items}, nil
}

func (s *OnboardingService) SelectWorkPlatform(userID uuid.UUID, req model.SelectWorkPlatformRequest) (*model.SelectWorkPlatformResponse, error) {
	platformID, err := uuid.Parse(req.WorkPlatformID)
	if err != nil {
		return nil, apperrors.BadRequest("work_platform_id tidak valid")
	}

	_, err = s.workPlatformRepo.GetByID(s.db, platformID)
	if err != nil {
		return nil, apperrors.NotFound("pekerjaan tidak ditemukan")
	}

	user, err := s.userRepo.GetUser(s.db, model.UserParam{UserID: userID})
	if err != nil {
		return nil, apperrors.NotFound("user tidak ditemukan")
	}

	user.WorkPlatformID = platformID
	if err := s.userRepo.UpdateUser(s.db, user); err != nil {
		return nil, apperrors.InternalServer("gagal menyimpan pekerjaan")
	}

	return &model.SelectWorkPlatformResponse{Message: "pekerjaan berhasil disimpan"}, nil
}

func (s *OnboardingService) GetIncomeSources() (*model.GetIncomeSourcesResponse, error) {
	sources := []model.IncomeSourceItem{
		{
			Type:        "manual",
			Label:       "Input Manual",
			Description: "Catat sendiri dalam kurang dari 10 detik. Tersedia di semua kondisi, tanpa koneksi platform.",
			IsAvailable: true,
		},
		{
			Type:        "ewallet",
			Label:       "Hubungkan GoPay / OVO",
			Description: "Transaksi masuk otomatis. Kami hanya membaca data masuk — tidak bisa transfer atau tarik saldo.",
			IsAvailable: false,
		},
	}
	return &model.GetIncomeSourcesResponse{Sources: sources}, nil
}

func (s *OnboardingService) SelectIncomeSource(userID uuid.UUID, req model.SelectIncomeSourceRequest) (*model.SelectIncomeSourceResponse, error) {
	user, err := s.userRepo.GetUser(s.db, model.UserParam{UserID: userID})
	if err != nil {
		return nil, apperrors.NotFound("user tidak ditemukan")
	}

	user.IncomeSourceType = req.IncomeSourceType
	if err := s.userRepo.UpdateUser(s.db, user); err != nil {
		return nil, apperrors.InternalServer("gagal menyimpan sumber penghasilan")
	}

	return &model.SelectIncomeSourceResponse{Message: "sumber penghasilan berhasil disimpan"}, nil
}
