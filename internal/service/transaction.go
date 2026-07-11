package service

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/azmiagr/sakutera-softdev/internal/repository"
	"github.com/azmiagr/sakutera-softdev/model"
	"github.com/azmiagr/sakutera-softdev/pkg/database/mariadb"
	apperr "github.com/azmiagr/sakutera-softdev/pkg/errors"
	"github.com/azmiagr/sakutera-softdev/pkg/mlclient"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITransactionService interface {
	GetSources(provider string) (*model.GetTransactionSourcesResponse, error)
	CreateTransaction(userID uuid.UUID, req model.CreateTransactionRequest) (*model.CreateTransactionResponse, error)
	GetTransactions(userID uuid.UUID, limit int) (*model.GetTransactionsResponse, error)
	GetLedger(userID uuid.UUID, period string, sourceID *uuid.UUID) (*model.GetLedgerResponse, error)
}

type TransactionService struct {
	db               *gorm.DB
	transactionRepo  repository.ITransactionRepository
	sourceRepo       repository.ITransactionSourceRepository
	forecastRepo     repository.IForecastResultRepository
	mlClient         mlclient.Interface
}

func NewTransactionService(
	transactionRepo repository.ITransactionRepository,
	sourceRepo repository.ITransactionSourceRepository,
	forecastRepo repository.IForecastResultRepository,
	mlClient mlclient.Interface,
) ITransactionService {
	return &TransactionService{
		db:              mariadb.Connection,
		transactionRepo: transactionRepo,
		sourceRepo:      sourceRepo,
		forecastRepo:    forecastRepo,
		mlClient:        mlClient,
	}
}

func (s *TransactionService) GetSources(provider string) (*model.GetTransactionSourcesResponse, error) {
	sources, err := s.sourceRepo.GetAll(s.db, provider)
	if err != nil {
		return nil, apperr.InternalServer("gagal mengambil sumber penghasilan")
	}

	var items []model.TransactionSourceItem
	for _, src := range sources {
		items = append(items, model.TransactionSourceItem{
			TransactionSourceID: src.TransactionSourceID.String(),
			Name:                src.Name,
			Provider:            src.Provider,
		})
	}

	return &model.GetTransactionSourcesResponse{Sources: items}, nil
}

func (s *TransactionService) CreateTransaction(userID uuid.UUID, req model.CreateTransactionRequest) (*model.CreateTransactionResponse, error) {
	sourceID, err := uuid.Parse(req.TransactionSourceID)
	if err != nil {
		return nil, apperr.BadRequest("transaction_source_id tidak valid")
	}

	source, err := s.sourceRepo.GetByID(s.db, sourceID)
	if err != nil {
		return nil, apperr.NotFound("sumber penghasilan tidak ditemukan")
	}

	txDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		return nil, apperr.BadRequest("format transaction_date tidak valid, gunakan YYYY-MM-DD")
	}

	prevHash := strings.Repeat("0", 64)
	lastTx, err := s.transactionRepo.GetLastByUserID(s.db, userID)
	if err == nil && lastTx != nil {
		prevHash = lastTx.CurrentHash
	}

	input := fmt.Sprintf("%s|%s|%.2f|%s|%s",
		userID.String(), sourceID.String(), req.Amount, req.TransactionDate, prevHash)
	currentHash := fmt.Sprintf("%x", sha256.Sum256([]byte(input)))

	t := &entity.Transaction{
		TransactionID:       uuid.New(),
		UserID:              userID,
		TransactionSourceID: sourceID,
		Amount:              req.Amount,
		TransactionDate:     txDate,
		Category:            req.Description,
		Status:              "success",
		PreviousHash:        prevHash,
		CurrentHash:         currentHash,
	}

	if err := s.transactionRepo.CreateTransaction(s.db, t); err != nil {
		return nil, apperr.InternalServer("gagal menyimpan transaksi")
	}

	go s.updateForecast(userID, source.Name)

	return &model.CreateTransactionResponse{
		TransactionID: t.TransactionID.String(),
		Amount:        t.Amount,
		CurrentHash:   t.CurrentHash,
		Status:        t.Status,
		Message:       "penghasilan berhasil dicatat ke ledger",
	}, nil
}

func (s *TransactionService) GetTransactions(userID uuid.UUID, limit int) (*model.GetTransactionsResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	transactions, err := s.transactionRepo.GetByUserID(s.db, userID, limit)
	if err != nil {
		return nil, apperr.InternalServer("gagal mengambil transaksi")
	}

	sourceIDs := make([]uuid.UUID, 0, len(transactions))
	for _, tx := range transactions {
		sourceIDs = append(sourceIDs, tx.TransactionSourceID)
	}
	sourceMap := s.fetchSourceMap(sourceIDs)

	var items []model.TransactionItem
	for _, tx := range transactions {
		src := sourceMap[tx.TransactionSourceID]
		hashPrefix := ""
		if len(tx.CurrentHash) >= 8 {
			hashPrefix = tx.CurrentHash[:8]
		}
		items = append(items, model.TransactionItem{
			TransactionID:   tx.TransactionID.String(),
			SourceName:      src.Name,
			SourceProvider:  src.Provider,
			Amount:          tx.Amount,
			TransactionDate: tx.TransactionDate.Format("2006-01-02"),
			Category:        tx.Category,
			HashPrefix:      hashPrefix,
		})
	}

	return &model.GetTransactionsResponse{Transactions: items, Total: len(items)}, nil
}

func (s *TransactionService) GetLedger(userID uuid.UUID, period string, sourceID *uuid.UUID) (*model.GetLedgerResponse, error) {
	var from, to *time.Time

	switch period {
	case "all", "":
		// no date filter
	case "bulan_ini":
		now := time.Now()
		f := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		t := now
		from, to = &f, &t
	case "3_bulan":
		f := time.Now().AddDate(0, -3, 0)
		t := time.Now()
		from, to = &f, &t
	default:
		return nil, apperr.BadRequest("period tidak valid, gunakan: all, bulan_ini, atau 3_bulan")
	}

	transactions, err := s.transactionRepo.GetLedger(s.db, userID, from, to, sourceID)
	if err != nil {
		return nil, apperr.InternalServer("gagal mengambil buku kas")
	}

	sourceIDs := make([]uuid.UUID, 0, len(transactions))
	for _, tx := range transactions {
		sourceIDs = append(sourceIDs, tx.TransactionSourceID)
	}
	sourceMap := s.fetchSourceMap(sourceIDs)

	var items []model.LedgerTransactionItem
	for _, tx := range transactions {
		src := sourceMap[tx.TransactionSourceID]
		hashPrefix := ""
		if len(tx.CurrentHash) >= 8 {
			hashPrefix = tx.CurrentHash[:8]
		}
		items = append(items, model.LedgerTransactionItem{
			TransactionID:   tx.TransactionID.String(),
			SourceName:      src.Name,
			SourceProvider:  src.Provider,
			Amount:          tx.Amount,
			TransactionDate: tx.TransactionDate.Format("2006-01-02"),
			TransactionTime: tx.TransactionDate.Format("15:04"),
			IsVerified:      true,
			HashPrefix:      hashPrefix,
		})
	}

	return &model.GetLedgerResponse{
		Summary: model.LedgerSummary{
			TotalEntries: int64(len(items)),
			ChainValid:   true,
		},
		Transactions: items,
		Total:        len(items),
	}, nil
}

func (s *TransactionService) fetchSourceMap(ids []uuid.UUID) map[uuid.UUID]*entity.TransactionSource {
	result := make(map[uuid.UUID]*entity.TransactionSource)
	for _, id := range ids {
		if _, ok := result[id]; !ok {
			src, err := s.sourceRepo.GetByID(s.db, id)
			if err == nil {
				result[id] = src
			} else {
				result[id] = &entity.TransactionSource{Name: "Lainnya", Provider: "Manual"}
			}
		}
	}
	return result
}

func (s *TransactionService) updateForecast(userID uuid.UUID, sourceName string) {
	from := time.Now().AddDate(0, 0, -30)
	to := time.Now()
	transactions, err := s.transactionRepo.GetByUserIDAndDateRange(s.db, userID, from, to)
	if err != nil || len(transactions) == 0 {
		return
	}

	var mlTxs []model.MLTransaction
	for _, tx := range transactions {
		mlTxs = append(mlTxs, model.MLTransaction{
			Date:   tx.TransactionDate.Format("2006-01-02"),
			Amount: int(tx.Amount),
			Source: strings.ToLower(sourceName),
		})
	}

	req := model.MLForecastRequest{
		UserID:       userID.String(),
		AsOfDate:     time.Now().Format("2006-01-02"),
		Transactions: mlTxs,
	}

	mlResp, err := s.mlClient.Forecast(req)
	if err != nil {
		return
	}

	forecastResult := &entity.ForecastResult{
		ForecastResultID:  uuid.New(),
		UserID:            userID,
		EMIValue:          mlResp.EMI.Value,
		Confidence:        mlResp.EMI.Confidence,
		TrendDirection:    mlResp.Trend.Direction,
		TrendChangePct:    mlResp.Trend.ChangePct,
		MonthToDateIncome: mlResp.MonthToDate.TotalEarned,
		ForecastTotal:     mlResp.ForecastRemainingMonth.ProjectedMonthTotal,
		RiskLevel:         mlResp.DeficitRisk.Level,
		RiskScore:         mlResp.DeficitRisk.Score,
		EstimatedShortfall: mlResp.DeficitRisk.EstimatedShortfall,
		ShouldNotify:      mlResp.DeficitRisk.ShouldNotify,
		IsDataSufficient:  mlResp.DataSufficiency.IsSufficient,
		DaysOfData:        mlResp.DataSufficiency.DaysOfData,
		TransactionCount:  len(transactions),
		ModelName:         "prophet",
	}

	_ = s.forecastRepo.UpsertByUserID(s.db, forecastResult)
}
