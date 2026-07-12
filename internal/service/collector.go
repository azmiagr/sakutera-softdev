package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/azmiagr/sakutera-softdev/internal/repository"
	"github.com/azmiagr/sakutera-softdev/model"
	"github.com/azmiagr/sakutera-softdev/pkg/database/mariadb"
	apperr "github.com/azmiagr/sakutera-softdev/pkg/errors"
	"github.com/azmiagr/sakutera-softdev/pkg/notifparser"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	maxBatchSize            = 50
	maxTitleLength          = 500
	maxTextLength           = 1000
	maxBigTextLength        = 2000
	maxPackageNameLength    = 255
	eventTimestampTolerance = 72 * time.Hour

	pairingCodeCharset       = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // tanpa 0/O/1/I biar gak ambigu
	pairingCodeLength        = 6
	defaultPairingCodeTTL    = 600 * time.Second
	minPairingCodeTTL        = 300 * time.Second
	maxPairingCodeTTL        = 900 * time.Second
	pairingRateLimitWindow   = 10 * time.Minute
	pairingRateLimitMaxCount = 5

	pairAttemptWindow   = 15 * time.Minute
	pairAttemptMaxCount = 10
)

type ICollectorService interface {
	CreatePairingCode(userID uuid.UUID, req model.CreatePairingCodeRequest) (*model.CreatePairingCodeResponse, error)
	PairDevice(req model.PairDeviceRequest, clientIP string) (*model.PairDeviceResponse, error)
	ListDevices(userID uuid.UUID) (*model.ListDevicesResponse, error)
	RevokeDevice(userID uuid.UUID, deviceID uuid.UUID) (*model.RevokeDeviceResponse, error)
	AuthenticateDevice(token string) (*entity.Device, error)
	HealthCheck(device *entity.Device) *model.CollectorHealthResponse
	IngestBatch(device *entity.Device, req model.NotificationBatchRequest) (*model.NotificationBatchResponse, error)
	GetConfig() (*model.CollectorConfigResponse, error)
	RedactOldNotifications(olderThan time.Duration) (int64, error)
}

type CollectorService struct {
	db                    *gorm.DB
	deviceRepo            repository.IDeviceRepository
	notificationEventRepo repository.INotificationEventRepository
	pairingCodeRepo       repository.IPairingCodeRepository
	transactionRepo       repository.ITransactionRepository
	transactionSourceRepo repository.ITransactionSourceRepository
	transactionReviewRepo repository.ITransactionReviewRepository
	collectorConfigRepo   repository.ICollectorConfigRepository
	auditLogRepo          repository.IAuditLogRepository
	rateLimitRepo         repository.IRateLimitRepository
}

func NewCollectorService(
	deviceRepo repository.IDeviceRepository,
	notificationEventRepo repository.INotificationEventRepository,
	pairingCodeRepo repository.IPairingCodeRepository,
	transactionRepo repository.ITransactionRepository,
	transactionSourceRepo repository.ITransactionSourceRepository,
	transactionReviewRepo repository.ITransactionReviewRepository,
	collectorConfigRepo repository.ICollectorConfigRepository,
	auditLogRepo repository.IAuditLogRepository,
	rateLimitRepo repository.IRateLimitRepository,
) ICollectorService {
	return &CollectorService{
		db:                    mariadb.Connection,
		deviceRepo:            deviceRepo,
		notificationEventRepo: notificationEventRepo,
		pairingCodeRepo:       pairingCodeRepo,
		transactionRepo:       transactionRepo,
		transactionSourceRepo: transactionSourceRepo,
		transactionReviewRepo: transactionReviewRepo,
		collectorConfigRepo:   collectorConfigRepo,
		auditLogRepo:          auditLogRepo,
		rateLimitRepo:         rateLimitRepo,
	}
}

func (s *CollectorService) writeAudit(userID, deviceID *uuid.UUID, action, detail string) {
	_ = s.auditLogRepo.Create(s.db, &entity.AuditLog{
		AuditLogID: uuid.New(),
		UserID:     userID,
		DeviceID:   deviceID,
		Action:     action,
		Detail:     detail,
	})
}

func generateDeviceToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func generatePairingCode() (string, error) {
	b := make([]byte, pairingCodeLength)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(pairingCodeCharset))))
		if err != nil {
			return "", err
		}
		b[i] = pairingCodeCharset[n.Int64()]
	}
	return fmt.Sprintf("SKT-%s", string(b)), nil
}

func (s *CollectorService) CreatePairingCode(userID uuid.UUID, req model.CreatePairingCodeRequest) (*model.CreatePairingCodeResponse, error) {
	count, err := s.pairingCodeRepo.CountRecentByUserID(s.db, userID, time.Now().Add(-pairingRateLimitWindow))
	if err != nil {
		return nil, apperr.InternalServer("gagal memeriksa rate limit")
	}
	if count >= pairingRateLimitMaxCount {
		return nil, apperr.TooManyRequests("terlalu banyak permintaan kode pairing, coba lagi nanti")
	}

	ttl := time.Duration(req.ExpiresInSeconds) * time.Second
	if ttl < minPairingCodeTTL || ttl > maxPairingCodeTTL {
		ttl = defaultPairingCodeTTL
	}

	if err := s.pairingCodeRepo.InvalidateActiveByUserID(s.db, userID); err != nil {
		return nil, apperr.InternalServer("gagal membatalkan kode pairing lama")
	}

	code, err := generatePairingCode()
	if err != nil {
		return nil, apperr.InternalServer("gagal membuat kode pairing")
	}

	pairingCode := &entity.PairingCode{
		PairingCodeID: uuid.New(),
		UserID:        userID,
		CodeHash:      hashToken(code),
		ExpiresAt:     time.Now().Add(ttl),
	}

	if err := s.pairingCodeRepo.Create(s.db, pairingCode); err != nil {
		return nil, apperr.InternalServer("gagal menyimpan kode pairing")
	}

	s.writeAudit(&userID, nil, "pairing_code.created", "")

	return &model.CreatePairingCodeResponse{
		PairingCode: code,
		ExpiresAt:   pairingCode.ExpiresAt.Format(time.RFC3339),
		QRPayload: model.QRPayload{
			Type:    "sakutera_device_pairing",
			Version: 1,
			Code:    code,
		},
	}, nil
}

func (s *CollectorService) PairDevice(req model.PairDeviceRequest, clientIP string) (*model.PairDeviceResponse, error) {
	rateLimitKey := "pair_attempt:" + clientIP
	count, err := s.rateLimitRepo.CountRecentByKey(s.db, rateLimitKey, time.Now().Add(-pairAttemptWindow))
	if err != nil {
		return nil, apperr.InternalServer("gagal memeriksa rate limit")
	}
	if count >= pairAttemptMaxCount {
		return nil, apperr.TooManyRequests("terlalu banyak percobaan pairing, coba lagi nanti")
	}
	_ = s.rateLimitRepo.Create(s.db, rateLimitKey)

	deviceID, err := uuid.Parse(req.DeviceID)
	if err != nil {
		return nil, apperr.BadRequest("device_id tidak valid")
	}

	token, err := generateDeviceToken()
	if err != nil {
		return nil, apperr.InternalServer("gagal membuat device token")
	}

	var resp *model.PairDeviceResponse

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		pairingCode, err := s.pairingCodeRepo.GetActiveByCodeHash(tx, hashToken(req.PairingCode))
		if err != nil {
			return apperr.Gone("pairing code sudah dipakai atau kedaluwarsa")
		}

		existing, err := s.deviceRepo.GetByID(tx, deviceID)
		if err == nil && existing.IsActive {
			return apperr.Conflict("perangkat sudah terhubung dan masih aktif")
		}

		pairedAt := time.Now()
		device := &entity.Device{
			DeviceID:   deviceID,
			UserID:     pairingCode.UserID,
			DeviceName: req.DeviceName,
			Platform:   req.Platform,
			OSVersion:  req.OSVersion,
			AppVersion: req.AppVersion,
			TokenHash:  hashToken(token),
			PairedAt:   pairedAt,
		}

		if err := s.deviceRepo.Upsert(tx, device); err != nil {
			return apperr.InternalServer("gagal menyimpan perangkat")
		}

		if err := s.pairingCodeRepo.MarkUsed(tx, pairingCode.PairingCodeID); err != nil {
			return apperr.InternalServer("gagal menandai kode pairing terpakai")
		}

		resp = &model.PairDeviceResponse{
			DeviceID:    deviceID.String(),
			DeviceToken: token,
			PairedAt:    pairedAt.Format(time.RFC3339),
		}

		return nil
	})

	if txErr != nil {
		var appErr *apperr.AppError
		if errors.As(txErr, &appErr) {
			return nil, appErr
		}
		return nil, apperr.InternalServer("gagal memproses pairing perangkat")
	}

	s.writeAudit(nil, &deviceID, "device.paired", "")

	return resp, nil
}

func (s *CollectorService) ListDevices(userID uuid.UUID) (*model.ListDevicesResponse, error) {
	devices, err := s.deviceRepo.GetAllByUserID(s.db, userID)
	if err != nil {
		return nil, apperr.InternalServer("gagal mengambil daftar perangkat")
	}

	items := make([]model.DeviceItem, 0, len(devices))
	for _, d := range devices {
		items = append(items, model.DeviceItem{
			DeviceID:   d.DeviceID.String(),
			DeviceName: d.DeviceName,
			Platform:   d.Platform,
			AppVersion: d.AppVersion,
			IsActive:   d.IsActive,
			PairedAt:   d.PairedAt.Format(time.RFC3339),
			LastSeenAt: d.LastSeenAt.Format(time.RFC3339),
		})
	}

	return &model.ListDevicesResponse{Devices: items}, nil
}

func (s *CollectorService) RevokeDevice(userID uuid.UUID, deviceID uuid.UUID) (*model.RevokeDeviceResponse, error) {
	device, err := s.deviceRepo.GetByID(s.db, deviceID)
	if err != nil || device.UserID != userID {
		return nil, apperr.NotFound("perangkat tidak ditemukan")
	}

	if err := s.deviceRepo.Revoke(s.db, deviceID); err != nil {
		return nil, apperr.InternalServer("gagal mencabut akses perangkat")
	}

	s.writeAudit(&userID, &deviceID, "device.revoked", "")

	return &model.RevokeDeviceResponse{
		DeviceID:  deviceID.String(),
		IsActive:  false,
		RevokedAt: time.Now().Format(time.RFC3339),
	}, nil
}

func (s *CollectorService) AuthenticateDevice(token string) (*entity.Device, error) {
	device, err := s.deviceRepo.GetByTokenHash(s.db, hashToken(token))
	if err != nil {
		return nil, apperr.Unauthorized("device token tidak valid")
	}

	if !device.IsActive || device.RevokedAt != nil {
		return nil, apperr.Unauthorized("device sudah tidak aktif")
	}

	_ = s.deviceRepo.UpdateLastSeen(s.db, device.DeviceID)

	return device, nil
}

func (s *CollectorService) HealthCheck(device *entity.Device) *model.CollectorHealthResponse {
	return &model.CollectorHealthResponse{
		Connected:    true,
		DeviceActive: device.IsActive,
		ServerTime:   time.Now().Format(time.RFC3339),
	}
}

func (s *CollectorService) IngestBatch(device *entity.Device, req model.NotificationBatchRequest) (*model.NotificationBatchResponse, error) {
	if req.DeviceID != device.DeviceID.String() {
		return nil, apperr.Forbidden("device_id tidak sesuai dengan token yang digunakan")
	}

	if len(req.Events) == 0 {
		return nil, apperr.BadRequest("events tidak boleh kosong")
	}

	if len(req.Events) > maxBatchSize {
		return nil, apperr.BadRequest("jumlah events melebihi batas maksimal")
	}

	resp := &model.NotificationBatchResponse{
		Received: len(req.Events),
		Results:  make([]model.NotificationEventResult, 0, len(req.Events)),
	}

	now := time.Now()

	for _, event := range req.Events {
		if errCode := validateEvent(event, now); errCode != "" {
			code := errCode
			resp.Failed++
			resp.Results = append(resp.Results, model.NotificationEventResult{
				EventID:   event.EventID,
				Status:    "failed",
				ErrorCode: &code,
			})
			continue
		}

		exists, err := s.notificationEventRepo.ExistsByEventID(s.db, event.EventID)
		if err != nil {
			code := "INTERNAL_ERROR"
			resp.Failed++
			resp.Results = append(resp.Results, model.NotificationEventResult{
				EventID:   event.EventID,
				Status:    "failed",
				ErrorCode: &code,
			})
			continue
		}

		if exists {
			resp.Duplicates++
			resp.Results = append(resp.Results, model.NotificationEventResult{
				EventID: event.EventID,
				Status:  "duplicate",
			})
			continue
		}

		postedAt := time.UnixMilli(event.PostedAt)

		notificationEvent := &entity.NotificationEvent{
			NotificationEventID: uuid.New(),
			EventID:             event.EventID,
			UserID:              device.UserID,
			DeviceID:            device.DeviceID,
			PackageName:         event.PackageName,
			NotificationID:      event.NotificationID,
			Title:               event.Title,
			Text:                event.Text,
			BigText:             event.BigText,
			PostedAt:            postedAt,
			CapturedAt:          time.UnixMilli(event.CapturedAt),
			ProcessingStatus:    "stored",
		}

		if err := s.notificationEventRepo.Create(s.db, notificationEvent); err != nil {
			code := "STORE_FAILED"
			resp.Failed++
			resp.Results = append(resp.Results, model.NotificationEventResult{
				EventID:   event.EventID,
				Status:    "failed",
				ErrorCode: &code,
			})
			continue
		}

		result := s.processNotificationEvent(device.UserID, notificationEvent, event, postedAt)
		switch result.Status {
		case "processed":
			resp.Processed++
		case "needs_review":
			resp.NeedsReview++
		case "ignored":
			resp.Ignored++
		}
		resp.Results = append(resp.Results, result)
	}

	s.writeAudit(&device.UserID, &device.DeviceID, "notifications.batch_ingested",
		fmt.Sprintf(`{"received":%d,"processed":%d,"needs_review":%d,"duplicates":%d,"ignored":%d,"failed":%d}`,
			resp.Received, resp.Processed, resp.NeedsReview, resp.Duplicates, resp.Ignored, resp.Failed))

	return resp, nil
}

// processNotificationEvent runs the parser against a stored raw event and
// decides whether it becomes a ledger transaction, a review candidate, or
// gets ignored. It always finalizes NotificationEvent.ProcessingStatus.
func (s *CollectorService) processNotificationEvent(
	userID uuid.UUID,
	notificationEvent *entity.NotificationEvent,
	event model.NotificationEventItem,
	postedAt time.Time,
) model.NotificationEventResult {
	allowed, err := s.collectorConfigRepo.IsPackageAllowed(s.db, event.PackageName)
	if err != nil || !allowed {
		_ = s.notificationEventRepo.UpdateProcessingStatus(s.db, notificationEvent.NotificationEventID, "ignored", notifparser.Version)
		code := "PACKAGE_NOT_ALLOWED"
		return model.NotificationEventResult{EventID: event.EventID, Status: "ignored", ErrorCode: &code}
	}

	candidate := notifparser.Parse(notifparser.Event{
		PackageName: event.PackageName,
		Title:       event.Title,
		Text:        event.Text,
		BigText:     event.BigText,
		PostedAt:    postedAt,
	})

	if candidate.Confidence == 0 || candidate.TransactionType == "expense" {
		_ = s.notificationEventRepo.UpdateProcessingStatus(s.db, notificationEvent.NotificationEventID, "ignored", notifparser.Version)
		return model.NotificationEventResult{EventID: event.EventID, Status: "ignored"}
	}

	source, err := s.transactionSourceRepo.GetByName(s.db, candidate.SourceName)
	if err == nil && candidate.Confidence >= notifparser.ConfidenceThreshold {
		notifEventID := notificationEvent.NotificationEventID
		t, _, err := createLedgerEntry(s.db, s.transactionRepo, s.transactionSourceRepo, userID, source.TransactionSourceID, candidate.Amount, candidate.TransactionDate, candidate.Description, &notifEventID)
		if err == nil {
			_ = s.notificationEventRepo.UpdateProcessingStatus(s.db, notificationEvent.NotificationEventID, "processed", notifparser.Version)
			txID := t.TransactionID.String()
			return model.NotificationEventResult{EventID: event.EventID, Status: "processed", TransactionID: &txID}
		}
	}

	var sourceID *uuid.UUID
	if err == nil {
		sourceID = &source.TransactionSourceID
	}

	review := &entity.TransactionReview{
		ReviewID:            uuid.New(),
		NotificationEventID: notificationEvent.NotificationEventID,
		UserID:              userID,
		Provider:            candidate.Provider,
		TransactionType:     candidate.TransactionType,
		Amount:              candidate.Amount,
		Description:         candidate.Description,
		TransactionDate:     candidate.TransactionDate,
		TransactionSourceID: sourceID,
		Confidence:          candidate.Confidence,
		Reason:              candidate.Reason,
	}

	if err := s.transactionReviewRepo.Create(s.db, review); err != nil {
		_ = s.notificationEventRepo.UpdateProcessingStatus(s.db, notificationEvent.NotificationEventID, "ignored", notifparser.Version)
		return model.NotificationEventResult{EventID: event.EventID, Status: "ignored"}
	}

	_ = s.notificationEventRepo.UpdateProcessingStatus(s.db, notificationEvent.NotificationEventID, "needs_review", notifparser.Version)
	reviewID := review.ReviewID.String()
	return model.NotificationEventResult{EventID: event.EventID, Status: "needs_review", ReviewID: &reviewID}
}

func validateEvent(event model.NotificationEventItem, now time.Time) string {
	if len(event.EventID) != 64 {
		return "INVALID_EVENT_ID"
	}
	if _, err := hex.DecodeString(event.EventID); err != nil {
		return "INVALID_EVENT_ID"
	}
	if event.PackageName == "" || len(event.PackageName) > maxPackageNameLength {
		return "INVALID_PACKAGE_NAME"
	}
	if len(event.Title) > maxTitleLength || len(event.Text) > maxTextLength || len(event.BigText) > maxBigTextLength {
		return "PAYLOAD_TOO_LARGE"
	}

	postedAt := time.UnixMilli(event.PostedAt)
	if postedAt.Before(now.Add(-eventTimestampTolerance)) || postedAt.After(now.Add(eventTimestampTolerance)) {
		return "TIMESTAMP_OUT_OF_RANGE"
	}

	return ""
}

func (s *CollectorService) GetConfig() (*model.CollectorConfigResponse, error) {
	config, err := s.collectorConfigRepo.GetConfig(s.db)
	if err != nil {
		return nil, apperr.InternalServer("gagal mengambil konfigurasi collector")
	}

	packages, err := s.collectorConfigRepo.GetEnabledPackages(s.db)
	if err != nil {
		return nil, apperr.InternalServer("gagal mengambil allowlist package")
	}

	items := make([]model.AllowedPackageItem, 0, len(packages))
	for _, p := range packages {
		items = append(items, model.AllowedPackageItem{
			PackageName: p.PackageName,
			Provider:    p.Provider,
			Enabled:     p.Enabled,
		})
	}

	return &model.CollectorConfigResponse{
		Mode:                config.Mode,
		ConfigVersion:       config.ConfigVersion,
		AllowedPackages:     items,
		MaxBatchSize:        maxBatchSize,
		SyncIntervalMinutes: 15,
	}, nil
}

func (s *CollectorService) RedactOldNotifications(olderThan time.Duration) (int64, error) {
	count, err := s.notificationEventRepo.RedactOldEvents(s.db, time.Now().Add(-olderThan))
	if err != nil {
		return 0, apperr.InternalServer("gagal redact notifikasi lama")
	}
	return count, nil
}
