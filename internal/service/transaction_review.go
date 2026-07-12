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

type ITransactionReviewService interface {
	GetReviews(userID uuid.UUID, status string, limit int) (*model.GetReviewsResponse, error)
	ConfirmReview(userID uuid.UUID, reviewID uuid.UUID, req model.ConfirmReviewRequest) (*model.ConfirmReviewResponse, error)
	RejectReview(userID uuid.UUID, reviewID uuid.UUID, req model.RejectReviewRequest) (*model.RejectReviewResponse, error)
}

type TransactionReviewService struct {
	db                    *gorm.DB
	transactionReviewRepo repository.ITransactionReviewRepository
	transactionRepo       repository.ITransactionRepository
	transactionSourceRepo repository.ITransactionSourceRepository
	auditLogRepo          repository.IAuditLogRepository
}

func NewTransactionReviewService(
	transactionReviewRepo repository.ITransactionReviewRepository,
	transactionRepo repository.ITransactionRepository,
	transactionSourceRepo repository.ITransactionSourceRepository,
	auditLogRepo repository.IAuditLogRepository,
) ITransactionReviewService {
	return &TransactionReviewService{
		db:                    mariadb.Connection,
		transactionReviewRepo: transactionReviewRepo,
		transactionRepo:       transactionRepo,
		transactionSourceRepo: transactionSourceRepo,
		auditLogRepo:          auditLogRepo,
	}
}

func (s *TransactionReviewService) writeAudit(userID uuid.UUID, action string) {
	_ = s.auditLogRepo.Create(s.db, &entity.AuditLog{
		AuditLogID: uuid.New(),
		UserID:     &userID,
		Action:     action,
	})
}

func (s *TransactionReviewService) GetReviews(userID uuid.UUID, status string, limit int) (*model.GetReviewsResponse, error) {
	if status == "" {
		status = "pending"
	}
	if limit <= 0 {
		limit = 20
	}

	reviews, err := s.transactionReviewRepo.GetByUserID(s.db, userID, status, limit)
	if err != nil {
		return nil, apperr.InternalServer("gagal mengambil kandidat transaksi")
	}

	items := make([]model.ReviewItem, 0, len(reviews))
	for _, r := range reviews {
		var sourceID *string
		if r.TransactionSourceID != nil {
			s := r.TransactionSourceID.String()
			sourceID = &s
		}

		items = append(items, model.ReviewItem{
			ReviewID:            r.ReviewID.String(),
			Provider:            r.Provider,
			TransactionType:     r.TransactionType,
			Amount:              r.Amount,
			Description:         r.Description,
			TransactionDate:     r.TransactionDate.Format(time.RFC3339),
			TransactionSourceID: sourceID,
			Confidence:          r.Confidence,
			Reason:              r.Reason,
			CreatedAt:           r.CreatedAt.Format(time.RFC3339),
		})
	}

	return &model.GetReviewsResponse{Items: items, NextCursor: nil}, nil
}

func (s *TransactionReviewService) ConfirmReview(userID uuid.UUID, reviewID uuid.UUID, req model.ConfirmReviewRequest) (*model.ConfirmReviewResponse, error) {
	review, err := s.transactionReviewRepo.GetByID(s.db, reviewID)
	if err != nil || review.UserID != userID {
		return nil, apperr.NotFound("kandidat transaksi tidak ditemukan")
	}

	if review.Status == "confirmed" && review.TransactionID != nil {
		return &model.ConfirmReviewResponse{
			ReviewID:      review.ReviewID.String(),
			TransactionID: review.TransactionID.String(),
			Status:        "confirmed",
		}, nil
	}

	if review.Status != "pending" {
		return nil, apperr.Conflict("kandidat transaksi sudah diproses sebelumnya")
	}

	sourceID, err := uuid.Parse(req.TransactionSourceID)
	if err != nil {
		return nil, apperr.BadRequest("transaction_source_id tidak valid")
	}

	txDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		txDate = review.TransactionDate
	}

	amount := req.Amount
	if amount <= 0 {
		amount = review.Amount
	}

	description := req.Description
	if description == "" {
		description = review.Description
	}

	notifEventID := review.NotificationEventID
	t, _, err := createLedgerEntry(s.db, s.transactionRepo, s.transactionSourceRepo, userID, sourceID, amount, txDate, description, &notifEventID)
	if err != nil {
		return nil, apperr.NotFound(err.Error())
	}

	if err := s.transactionReviewRepo.UpdateStatus(s.db, reviewID, "confirmed", &t.TransactionID); err != nil {
		return nil, apperr.InternalServer("gagal memperbarui status kandidat transaksi")
	}

	s.writeAudit(userID, "review.confirmed")

	return &model.ConfirmReviewResponse{
		ReviewID:      review.ReviewID.String(),
		TransactionID: t.TransactionID.String(),
		Status:        "confirmed",
	}, nil
}

func (s *TransactionReviewService) RejectReview(userID uuid.UUID, reviewID uuid.UUID, req model.RejectReviewRequest) (*model.RejectReviewResponse, error) {
	review, err := s.transactionReviewRepo.GetByID(s.db, reviewID)
	if err != nil || review.UserID != userID {
		return nil, apperr.NotFound("kandidat transaksi tidak ditemukan")
	}

	if review.Status != "pending" {
		return nil, apperr.Conflict("kandidat transaksi sudah diproses sebelumnya")
	}

	if err := s.transactionReviewRepo.UpdateStatus(s.db, reviewID, "rejected", nil); err != nil {
		return nil, apperr.InternalServer("gagal menolak kandidat transaksi")
	}

	s.writeAudit(userID, "review.rejected")

	return &model.RejectReviewResponse{
		ReviewID: review.ReviewID.String(),
		Status:   "rejected",
	}, nil
}
