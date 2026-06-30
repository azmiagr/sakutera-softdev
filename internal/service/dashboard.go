package service

import (
	"time"

	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/azmiagr/sakutera-softdev/internal/repository"
	"github.com/azmiagr/sakutera-softdev/model"
	"github.com/azmiagr/sakutera-softdev/pkg/database/mariadb"
	apperr "github.com/azmiagr/sakutera-softdev/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IDashboardService interface {
	GetDashboard(userID uuid.UUID) (*model.DashboardResponse, error)
}

type DashboardService struct {
	db              *gorm.DB
	userRepo        repository.IUserRepository
	transactionRepo repository.ITransactionRepository
	sourceRepo      repository.ITransactionSourceRepository
	forecastRepo    repository.IForecastResultRepository
}

func NewDashboardService(
	userRepo repository.IUserRepository,
	transactionRepo repository.ITransactionRepository,
	sourceRepo repository.ITransactionSourceRepository,
	forecastRepo repository.IForecastResultRepository,
) IDashboardService {
	return &DashboardService{
		db:              mariadb.Connection,
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
		sourceRepo:      sourceRepo,
		forecastRepo:    forecastRepo,
	}
}

func (s *DashboardService) GetDashboard(userID uuid.UUID) (*model.DashboardResponse, error) {
	user, err := s.userRepo.GetUser(s.db, model.UserParam{UserID: userID})
	if err != nil {
		return nil, apperr.NotFound("user tidak ditemukan")
	}

	var forecastData *model.ForecastData
	if forecast, err := s.forecastRepo.GetByUserID(s.db, userID); err == nil {
		forecastData = &model.ForecastData{
			EMIValue:          forecast.EMIValue,
			Confidence:        forecast.Confidence,
			TrendDirection:    forecast.TrendDirection,
			TrendChangePct:    forecast.TrendChangePct,
			MonthToDateIncome: forecast.MonthToDateIncome,
			ForecastTotal:     forecast.ForecastTotal,
			RiskLevel:         forecast.RiskLevel,
			RiskScore:         forecast.RiskScore,
			IsDataSufficient:  forecast.IsDataSufficient,
			DaysOfData:        forecast.DaysOfData,
			TransactionCount:  forecast.TransactionCount,
			ModelName:         forecast.ModelName,
		}
	}

	todayIncome, _ := s.transactionRepo.SumTodayByUserID(s.db, userID)
	validChainCount, _ := s.transactionRepo.CountSuccessByUserID(s.db, userID)

	recentTxs, _ := s.transactionRepo.GetByUserID(s.db, userID, 5)
	recentItems := s.buildTransactionItems(recentTxs)

	from := time.Now().AddDate(0, 0, -29)
	to := time.Now()
	rangeTxs, _ := s.transactionRepo.GetByUserIDAndDateRange(s.db, userID, from, to)
	trendData := buildDailyTrend(rangeTxs, from)

	return &model.DashboardResponse{
		UserFullName:       user.FullName,
		Forecast:           forecastData,
		TodayIncome:        todayIncome,
		ValidChainCount:    validChainCount,
		RecentTransactions: recentItems,
		TrendData:          trendData,
	}, nil
}

func (s *DashboardService) buildTransactionItems(txs []entity.Transaction) []model.TransactionItem {
	var items []model.TransactionItem
	for _, tx := range txs {
		srcName, srcProvider := "Lainnya", "Manual"
		if src, err := s.sourceRepo.GetByID(s.db, tx.TransactionSourceID); err == nil {
			srcName = src.Name
			srcProvider = src.Provider
		}

		hashPrefix := ""
		if len(tx.CurrentHash) >= 8 {
			hashPrefix = tx.CurrentHash[:8]
		}

		items = append(items, model.TransactionItem{
			TransactionID:   tx.TransactionID.String(),
			SourceName:      srcName,
			SourceProvider:  srcProvider,
			Amount:          tx.Amount,
			TransactionDate: tx.TransactionDate.Format("2006-01-02"),
			Category:        tx.Category,
			HashPrefix:      hashPrefix,
		})
	}
	return items
}

func buildDailyTrend(txs []entity.Transaction, from time.Time) []model.TrendPoint {
	dailyMap := make(map[string]float64)
	for _, tx := range txs {
		day := tx.TransactionDate.Format("2006-01-02")
		dailyMap[day] += tx.Amount
	}

	var points []model.TrendPoint
	for i := 0; i < 30; i++ {
		day := from.AddDate(0, 0, i).Format("2006-01-02")
		points = append(points, model.TrendPoint{
			Date:   day,
			Amount: dailyMap[day],
		})
	}
	return points
}
