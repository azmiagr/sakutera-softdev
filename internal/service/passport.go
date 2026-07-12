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

var periodOrder = []string{"3_bulan", "6_bulan", "12_bulan"}

type periodRequirement struct {
	label    string
	minDays  int // reduced eligibility threshold (passport can be issued once this much data exists)
	fullDays int // true full-period length in days, used for is_full_period and EMI normalization
}

var periodRequirements = map[string]periodRequirement{
	"3_bulan":  {label: "3 Bulan", minDays: 30, fullDays: 90},
	"6_bulan":  {label: "6 Bulan", minDays: 180, fullDays: 180},
	"12_bulan": {label: "12 Bulan", minDays: 365, fullDays: 365},
}

type IPassportService interface {
	GetPassport(userID uuid.UUID) (*model.GetPassportResponse, error)
	PreviewPassport(userID uuid.UUID, period string) (*model.PassportPreviewResponse, error)
	IssuePassport(userID uuid.UUID, period string) (*model.IssuePassportResponse, error)
}

type PassportService struct {
	db              *gorm.DB
	passportRepo    repository.IIncomePassportRepository
	transactionRepo repository.ITransactionRepository
	forecastRepo    repository.IForecastResultRepository
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
	req, ok := periodRequirements[period]
	if !ok {
		return nil, apperr.BadRequest("period tidak valid, gunakan: 3_bulan, 6_bulan, atau 12_bulan")
	}

	daysAvailable := s.daysOfData(userID)
	if daysAvailable < req.minDays {
		return nil, apperr.BadRequestWithData(
			fmt.Sprintf("data belum cukup untuk passport periode %s", req.label),
			model.EligibilityGapDetail{
				Period:        period,
				DaysAvailable: daysAvailable,
				DaysRequired:  req.minDays,
				RemainingDays: req.minDays - daysAvailable,
			},
		)
	}

	forecast, err := s.forecastRepo.GetByUserID(s.db, userID)
	if err != nil {
		return nil, apperr.BadRequest("data forecast belum tersedia, catat minimal 30 hari transaksi terlebih dahulu")
	}

	daysUsed, from, to := periodWindow(daysAvailable, req.fullDays)

	txs, err := s.transactionRepo.GetByUserIDAndDateRange(s.db, userID, from, to)
	if err != nil {
		return nil, apperr.InternalServer("gagal mengambil data transaksi")
	}

	emiValue := sumTransactions(txs) / (float64(daysUsed) / 30.0)

	return &model.PassportPreviewResponse{
		PeriodType:     period,
		PeriodLabel:    buildPeriodLabel(period, from, to, false),
		PeriodStart:    from.Format("2006-01-02"),
		PeriodEnd:      to.Format("2006-01-02"),
		DaysUsed:       daysUsed,
		DaysRequired:   req.minDays,
		IsFullPeriod:   daysAvailable >= req.fullDays,
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
	req, ok := periodRequirements[period]
	if !ok {
		return nil, apperr.BadRequest("period tidak valid, gunakan: 3_bulan, 6_bulan, atau 12_bulan")
	}

	daysAvailable := s.daysOfData(userID)
	if daysAvailable < req.minDays {
		return nil, apperr.BadRequestWithData(
			fmt.Sprintf("data belum cukup untuk menerbitkan passport periode %s", req.label),
			model.EligibilityGapDetail{
				DaysAvailable: daysAvailable,
				DaysRequired:  req.minDays,
				RemainingDays: req.minDays - daysAvailable,
			},
		)
	}

	forecast, err := s.forecastRepo.GetByUserID(s.db, userID)
	if err != nil {
		return nil, apperr.BadRequest("data forecast belum tersedia, catat minimal 30 hari transaksi terlebih dahulu")
	}

	daysUsed, from, to := periodWindow(daysAvailable, req.fullDays)

	txs, err := s.transactionRepo.GetByUserIDAndDateRange(s.db, userID, from, to)
	if err != nil {
		return nil, apperr.InternalServer("gagal mengambil data transaksi")
	}

	emiValue := sumTransactions(txs) / (float64(daysUsed) / 30.0)

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

// buildEligibility constructs per-period eligibility data, using ForecastResult if available.
func (s *PassportService) buildEligibility(userID uuid.UUID) model.PassportEligibility {
	daysOfData, entriesVerified := s.daysOfDataAndEntries(userID)

	periods := make([]model.PeriodEligibility, 0, len(periodOrder))
	anyEligible := false
	for _, key := range periodOrder {
		req := periodRequirements[key]
		eligible := daysOfData >= req.minDays
		if eligible {
			anyEligible = true
		}

		remaining := max(req.minDays-daysOfData, 0)

		periods = append(periods, model.PeriodEligibility{
			Period:        key,
			Label:         req.label,
			IsEligible:    eligible,
			DaysRequired:  req.minDays,
			DaysAvailable: daysOfData,
			RemainingDays: remaining,
		})
	}

	return model.PassportEligibility{
		IsEligible:      anyEligible,
		DaysOfData:      daysOfData,
		EntriesVerified: entriesVerified,
		Periods:         periods,
	}
}

func (s *PassportService) daysOfData(userID uuid.UUID) int {
	days, _ := s.daysOfDataAndEntries(userID)
	return days
}

func (s *PassportService) daysOfDataAndEntries(userID uuid.UUID) (int, int64) {
	if forecast, err := s.forecastRepo.GetByUserID(s.db, userID); err == nil {
		return forecast.DaysOfData, int64(forecast.TransactionCount)
	}

	count, _ := s.transactionRepo.CountSuccessByUserID(s.db, userID)
	return int(count), count
}

// periodWindow returns (daysUsed, periodStart, periodEnd) reflecting the
// actual data coverage available: capped at fullDays, ending today.
func periodWindow(daysAvailable, fullDays int) (int, time.Time, time.Time) {
	daysUsed := min(daysAvailable, fullDays)

	to := time.Now()
	from := to.AddDate(0, 0, -daysUsed)
	return daysUsed, from, to
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
