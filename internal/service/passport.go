package service

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/azmiagr/sakutera-softdev/internal/repository"
	"github.com/azmiagr/sakutera-softdev/model"
	"github.com/azmiagr/sakutera-softdev/pkg/database/mariadb"
	apperr "github.com/azmiagr/sakutera-softdev/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const minDaysRequired = 30

type IPassportService interface {
	GetPassport(userID uuid.UUID) (*model.GetPassportResponse, error)
	PreviewPassport(userID uuid.UUID, period string) (*model.PassportPreviewResponse, error)
	IssuePassport(userID uuid.UUID, period string) (*model.IssuePassportResponse, error)
}

type PassportService struct {
	db             *gorm.DB
	passportRepo   repository.IIncomePassportRepository
	transactionRepo repository.ITransactionRepository
	forecastRepo   repository.IForecastResultRepository
}

func NewPassportService(
	passportRepo repository.IIncomePassportRepository,
	transactionRepo repository.ITransactionRepository,
	forecastRepo repository.IForecastResultRepository,
) IPassportService {
	return &PassportService{
		db:              mariadb.Connection,
		passportRepo:    passportRepo,
		transactionRepo: transactionRepo,
		forecastRepo:    forecastRepo,
	}
}

func (s *PassportService) GetPassport(userID uuid.UUID) (*model.GetPassportResponse, error) {
	eligibility := s.buildEligibility(userID)

	var activePassport *model.ActivePassportItem
	p, err := s.passportRepo.GetLatestByUserID(s.db, userID)
	if err == nil {
		activePassport = &model.ActivePassportItem{
			IncomePassportID: p.IncomePassportID.String(),
			PassportNumber:   p.PassportNumber,
			EMIValue:         p.EMI,
			PeriodType:       p.PeriodType,
			PeriodLabel:      buildPeriodLabel(p.PeriodType, p.PeriodStart, p.PeriodEnd, true),
			RiskLevel:        p.RiskLevel,
			IssuedAt:         p.CreatedAt.Format("2006-01-02"),
		}
	}

	return &model.GetPassportResponse{
		Eligibility:    eligibility,
		ActivePassport: activePassport,
	}, nil
}

func (s *PassportService) PreviewPassport(userID uuid.UUID, period string) (*model.PassportPreviewResponse, error) {
	numMonths, err := periodToMonths(period)
	if err != nil {
		return nil, err
	}

	from, to := calcDateRange(numMonths)

	forecast, err := s.forecastRepo.GetByUserID(s.db, userID)
	if err != nil {
		return nil, apperr.BadRequest("data forecast belum tersedia, catat minimal 30 hari transaksi terlebih dahulu")
	}

	txs, err := s.transactionRepo.GetByUserIDAndDateRange(s.db, userID, from, to)
	if err != nil {
		return nil, apperr.InternalServer("gagal mengambil data transaksi")
	}

	emiValue := sumTransactions(txs) / float64(numMonths)

	return &model.PassportPreviewResponse{
		PeriodType:     period,
		PeriodLabel:    buildPeriodLabel(period, from, to, false),
		PeriodStart:    from.Format("2006-01-02"),
		PeriodEnd:      to.Format("2006-01-02"),
		EMIValue:       emiValue,
		StabilityLabel: stabilityLabel(forecast.TrendDirection),
		TrendDirection: forecast.TrendDirection,
		TrendChangePct: forecast.TrendChangePct,
		RiskLevel:      forecast.RiskLevel,
		RiskScore:      forecast.RiskScore,
		TotalEntries:   len(txs),
	}, nil
}

func (s *PassportService) IssuePassport(userID uuid.UUID, period string) (*model.IssuePassportResponse, error) {
	numMonths, err := periodToMonths(period)
	if err != nil {
		return nil, err
	}

	forecast, err := s.forecastRepo.GetByUserID(s.db, userID)
	if err != nil {
		return nil, apperr.BadRequest("data forecast belum tersedia, catat minimal 30 hari transaksi terlebih dahulu")
	}

	from, to := calcDateRange(numMonths)

	txs, err := s.transactionRepo.GetByUserIDAndDateRange(s.db, userID, from, to)
	if err != nil {
		return nil, apperr.InternalServer("gagal mengambil data transaksi")
	}

	emiValue := sumTransactions(txs) / float64(numMonths)

	lastTx, err := s.transactionRepo.GetLastByUserID(s.db, userID)
	if err != nil {
		return nil, apperr.BadRequest("belum ada transaksi yang tercatat")
	}

	passportNumber := fmt.Sprintf("SKT-%d-%s-%s",
		time.Now().Year(),
		randomTwoUpper(),
		strings.ToUpper(lastTx.CurrentHash[:8]),
	)

	p := &entity.IncomePassport{
		IncomePassportID: uuid.New(),
		UserID:           userID,
		PassportNumber:   passportNumber,
		PeriodType:       period,
		PeriodStart:      from,
		PeriodEnd:        to,
		EMI:              emiValue,
		StabilityScore:   1.0 - forecast.RiskScore,
		StabilityLabel:   stabilityLabel(forecast.TrendDirection),
		TrendDirection:   forecast.TrendDirection,
		TrendChangePct:   forecast.TrendChangePct,
		RiskLevel:        forecast.RiskLevel,
		RiskScore:        forecast.RiskScore,
		TotalEntries:     len(txs),
	}

	if err := s.passportRepo.Create(s.db, p); err != nil {
		return nil, apperr.InternalServer("gagal menerbitkan income passport")
	}

	return &model.IssuePassportResponse{
		IncomePassportID: p.IncomePassportID.String(),
		PassportNumber:   p.PassportNumber,
		EMIValue:         p.EMI,
		PeriodType:       p.PeriodType,
		PeriodLabel:      buildPeriodLabel(period, from, to, true),
		RiskLevel:        p.RiskLevel,
		IssuedAt:         p.CreatedAt.Format("2006-01-02"),
	}, nil
}

// buildEligibility constructs eligibility data, using ForecastResult if available.
func (s *PassportService) buildEligibility(userID uuid.UUID) model.PassportEligibility {
	daysOfData := 0
	var entriesVerified int64

	if forecast, err := s.forecastRepo.GetByUserID(s.db, userID); err == nil {
		daysOfData = forecast.DaysOfData
		entriesVerified = int64(forecast.TransactionCount)
	} else {
		count, _ := s.transactionRepo.CountSuccessByUserID(s.db, userID)
		entriesVerified = count
		daysOfData = int(count)
	}

	return model.PassportEligibility{
		IsEligible:      daysOfData >= minDaysRequired,
		DaysOfData:      daysOfData,
		MinRequired:     minDaysRequired,
		EntriesVerified: entriesVerified,
	}
}

func periodToMonths(period string) (int, error) {
	switch period {
	case "3_bulan":
		return 3, nil
	case "6_bulan":
		return 6, nil
	case "12_bulan":
		return 12, nil
	default:
		return 0, apperr.BadRequest("period tidak valid, gunakan: 3_bulan, 6_bulan, atau 12_bulan")
	}
}

// calcDateRange returns (firstDayOfStartMonth, lastDayOfCurrentMonth).
func calcDateRange(numMonths int) (time.Time, time.Time) {
	now := time.Now()
	firstCurrentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	from := firstCurrentMonth.AddDate(0, -(numMonths - 1), 0)
	to := firstCurrentMonth.AddDate(0, 1, -1)
	return from, to
}

// buildPeriodLabel builds display labels like "Apr–Jun 2026" (withYear=true) or "Apr–Jun" (withYear=false).
// For 12-month periods it uses 2-digit years: "Jun 25–Jun 26".
func buildPeriodLabel(period string, from, to time.Time, withYear bool) string {
	if period == "12_bulan" {
		return fmt.Sprintf("%s %02d–%s %02d",
			from.Format("Jan"), from.Year()%100,
			to.Format("Jan"), to.Year()%100,
		)
	}
	if withYear {
		return fmt.Sprintf("%s–%s %d", from.Format("Jan"), to.Format("Jan"), to.Year())
	}
	return fmt.Sprintf("%s–%s", from.Format("Jan"), to.Format("Jan"))
}

func stabilityLabel(trendDirection string) string {
	switch trendDirection {
	case "up":
		return "STABIL"
	case "stable":
		return "CUKUP STABIL"
	default:
		return "FLUKTUATIF"
	}
}

func sumTransactions(txs []entity.Transaction) float64 {
	var total float64
	for _, tx := range txs {
		total += tx.Amount
	}
	return total
}

func randomTwoUpper() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 2)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
