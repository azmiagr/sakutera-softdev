package service

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/azmiagr/sakutera-softdev/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func createLedgerEntry(
	db *gorm.DB,
	transactionRepo repository.ITransactionRepository,
	sourceRepo repository.ITransactionSourceRepository,
	userID, sourceID uuid.UUID,
	amount float64,
	transactionDate time.Time,
	category string,
	notificationEventID *uuid.UUID,
) (*entity.Transaction, *entity.TransactionSource, error) {
	source, err := sourceRepo.GetByID(db, sourceID)
	if err != nil {
		return nil, nil, fmt.Errorf("sumber penghasilan tidak ditemukan")
	}

	prevHash := strings.Repeat("0", 64)
	lastTx, err := transactionRepo.GetLastByUserID(db, userID)
	if err == nil && lastTx != nil {
		prevHash = lastTx.CurrentHash
	}

	dateStr := transactionDate.Format("2006-01-02")
	input := fmt.Sprintf("%s|%s|%.2f|%s|%s", userID.String(), sourceID.String(), amount, dateStr, prevHash)
	currentHash := fmt.Sprintf("%x", sha256.Sum256([]byte(input)))

	t := &entity.Transaction{
		TransactionID:             uuid.New(),
		UserID:                    userID,
		TransactionSourceID:       sourceID,
		Amount:                    amount,
		TransactionDate:           transactionDate,
		Category:                  category,
		Status:                    "success",
		PreviousHash:              prevHash,
		CurrentHash:               currentHash,
		SourceNotificationEventID: notificationEventID,
	}

	if err := transactionRepo.CreateTransaction(db, t); err != nil {
		return nil, nil, fmt.Errorf("gagal menyimpan transaksi")
	}

	return t, source, nil
}
