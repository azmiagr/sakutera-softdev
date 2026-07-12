package service

import (
	"fmt"
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

type IAccessService interface {
	GetConsents(userID uuid.UUID) (*model.GetConsentsResponse, error)
	GrantAccess(userID uuid.UUID, req model.GrantAccessRequest) error
	RevokeAccess(userID uuid.UUID, consentID string) error
	GetAccessLogs(userID uuid.UUID, filter string) (*model.GetAccessLogsResponse, error)
	GetOrganizations() ([]model.OrganizationItem, error)
}

type AccessService struct {
	db           *gorm.DB
	passportRepo repository.IIncomePassportRepository
	consentRepo  repository.IConsentRepository
	logRepo      repository.IAccessLogRepository
	orgRepo      repository.IOrganizationRepository
}

func NewAccessService(
	passportRepo repository.IIncomePassportRepository,
	consentRepo repository.IConsentRepository,
	logRepo repository.IAccessLogRepository,
	orgRepo repository.IOrganizationRepository,
) IAccessService {
	return &AccessService{
		db:           mariadb.Connection,
		passportRepo: passportRepo,
		consentRepo:  consentRepo,
		logRepo:      logRepo,
		orgRepo:      orgRepo,
	}
}

func (s *AccessService) GetConsents(userID uuid.UUID) (*model.GetConsentsResponse, error) {
	passports, err := s.passportRepo.GetAllByUserID(s.db, userID)
	if err != nil || len(passports) == 0 {
		return &model.GetConsentsResponse{Consents: []model.ConsentItem{}}, nil
	}

	passportMap, passportIDs := indexPassports(passports)

	consents, err := s.consentRepo.GetByPassportIDs(s.db, passportIDs)
	if err != nil {
		return nil, apperr.InternalServer("gagal mengambil data akses")
	}

	items := make([]model.ConsentItem, 0, len(consents))
	for _, c := range consents {
		org, err := s.orgRepo.GetByID(s.db, c.OrganizationID)
		if err != nil {
			continue
		}

		passport := passportMap[c.PassportID]

		item := model.ConsentItem{
			ConsentID:        c.ConsentID.String(),
			OrganizationName: org.Name,
			OrganizationType: org.Type,
			GrantedAt:        c.GrantedAt.Format("02 Jan 2006"),
			DataScope:        scopeToLabels(c.DataScope),
			Status:           c.Status,
			StatusLabel:      consentStatusLabel(c.Status, c.ExpiresAt),
			Purpose:          c.Purpose,
			IncomePassportID: passport.IncomePassportID.String(),
			PassportNumber:   passport.PassportNumber,
			PeriodType:       passport.PeriodType,
			PeriodLabel:      buildPeriodLabel(passport.PeriodType, passport.PeriodStart, passport.PeriodEnd, true),
		}

		if c.ExpiresAt != nil {
			formatted := c.ExpiresAt.Format("02 Jan 2006")
			item.ExpiresAt = &formatted
			days := max(int(time.Until(*c.ExpiresAt).Hours()/24), 0)
			item.DaysRemaining = &days
		}

		items = append(items, item)
	}

	return &model.GetConsentsResponse{Consents: items}, nil
}

func (s *AccessService) GrantAccess(userID uuid.UUID, req model.GrantAccessRequest) error {
	passportID, err := uuid.Parse(req.IncomePassportID)
	if err != nil {
		return apperr.BadRequest("income_passport_id tidak valid")
	}

	passport, err := s.passportRepo.GetByID(s.db, passportID)
	if err != nil || passport.UserID != userID {
		return apperr.NotFound("income passport tidak ditemukan")
	}

	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		return apperr.BadRequest("organization_id tidak valid")
	}

	org, err := s.orgRepo.GetByID(s.db, orgID)
	if err != nil {
		return apperr.NotFound("organisasi tidak ditemukan")
	}

	existing, _ := s.consentRepo.GetByPassportAndOrg(s.db, passport.IncomePassportID, orgID)
	if existing != nil {
		return apperr.Conflict("akses untuk organisasi ini sudah aktif")
	}

	now := time.Now()
	consent := &entity.Consent{
		ConsentID:      uuid.New(),
		PassportID:     passport.IncomePassportID,
		OrganizationID: orgID,
		Permission:     "read",
		Status:         "active",
		DataScope:      strings.Join(req.DataScope, ","),
		Purpose:        req.Purpose,
		GrantedAt:      now,
	}

	if req.ExpiresInDays > 0 {
		exp := now.AddDate(0, 0, req.ExpiresInDays)
		consent.ExpiresAt = &exp
	}

	err = s.consentRepo.Create(s.db, consent)
	if err != nil {
		return apperr.InternalServer("gagal memberikan akses")
	}

	scopeLabel := scopeToText(req.DataScope)
	note := fmt.Sprintf("Akses pertama · %s", scopeLabel)
	if req.Purpose != "" {
		note = fmt.Sprintf("Akses pertama · %s · %s", scopeLabel, req.Purpose)
	}

	logEntry := &entity.AccessLog{
		AccessLogID:    uuid.New(),
		PassportID:     passport.IncomePassportID,
		OrganizationID: orgID,
		AccessedAt:     now,
		Status:         "success",
		Note:           note,
	}
	_ = s.logRepo.Create(s.db, logEntry)

	_ = org
	return nil
}

func (s *AccessService) RevokeAccess(userID uuid.UUID, consentID string) error {
	cID, err := uuid.Parse(consentID)
	if err != nil {
		return apperr.BadRequest("consent_id tidak valid")
	}

	consent, err := s.consentRepo.GetByConsentID(s.db, cID)
	if err != nil {
		return apperr.NotFound("akses tidak ditemukan")
	}

	passport, err := s.passportRepo.GetByID(s.db, consent.PassportID)
	if err != nil || passport.UserID != userID {
		return apperr.Unauthorized("tidak memiliki izin untuk mencabut akses ini")
	}

	if consent.Status != "active" {
		return apperr.BadRequest("akses sudah tidak aktif")
	}

	err = s.consentRepo.UpdateStatus(s.db, cID, "revoked")
	if err != nil {
		return apperr.InternalServer("gagal mencabut akses")
	}

	now := time.Now()
	logEntry := &entity.AccessLog{
		AccessLogID:    uuid.New(),
		PassportID:     consent.PassportID,
		OrganizationID: consent.OrganizationID,
		AccessedAt:     now,
		Status:         "success",
		Note:           fmt.Sprintf("Akses dicabut oleh kamu · %s", now.Format("2 Jan")),
	}
	_ = s.logRepo.Create(s.db, logEntry)

	return nil
}

func (s *AccessService) GetAccessLogs(userID uuid.UUID, filter string) (*model.GetAccessLogsResponse, error) {
	passports, err := s.passportRepo.GetAllByUserID(s.db, userID)
	if err != nil || len(passports) == 0 {
		return &model.GetAccessLogsResponse{Logs: []model.AccessLogItem{}}, nil
	}

	passportMap, passportIDs := indexPassports(passports)

	logs, err := s.logRepo.GetByPassportIDs(s.db, passportIDs)
	if err != nil {
		return nil, apperr.InternalServer("gagal mengambil riwayat akses")
	}

	consentCache := map[string]*entity.Consent{}

	items := make([]model.AccessLogItem, 0, len(logs))
	for _, l := range logs {
		org, err := s.orgRepo.GetByID(s.db, l.OrganizationID)
		if err != nil {
			continue
		}

		passport := passportMap[l.PassportID]

		cacheKey := l.PassportID.String() + ":" + l.OrganizationID.String()
		consent, ok := consentCache[cacheKey]
		if !ok {
			c, _ := s.consentRepo.GetByPassportAndOrg(s.db, l.PassportID, l.OrganizationID)
			consentCache[cacheKey] = c
			consent = c
		}

		consentStatus := "revoked"
		statusLabel := "DICABUT"
		var dataScope []string

		if consent != nil {
			consentStatus = consent.Status
			dataScope = scopeToLabels(consent.DataScope)
			if consent.Status == "active" {
				statusLabel = "VALID"
			}
		}

		items = append(items, model.AccessLogItem{
			AccessLogID:      l.AccessLogID.String(),
			OrganizationName: org.Name,
			OrganizationType: org.Type,
			AccessedAt:       l.AccessedAt.Format("02 Jan 2006 · 15:04 WIB"),
			DataScope:        dataScope,
			ConsentStatus:    consentStatus,
			StatusLabel:      statusLabel,
			Note:             l.Note,
			IncomePassportID: passport.IncomePassportID.String(),
			PassportNumber:   passport.PassportNumber,
			PeriodType:       passport.PeriodType,
			PeriodLabel:      buildPeriodLabel(passport.PeriodType, passport.PeriodStart, passport.PeriodEnd, true),
		})
	}

	return &model.GetAccessLogsResponse{Logs: items}, nil
}

func (s *AccessService) GetOrganizations() ([]model.OrganizationItem, error) {
	orgs, err := s.orgRepo.GetAll(s.db)
	if err != nil {
		return nil, apperr.InternalServer("gagal mengambil daftar organisasi")
	}

	items := make([]model.OrganizationItem, 0, len(orgs))
	for _, o := range orgs {
		items = append(items, model.OrganizationItem{
			OrganizationID: o.OrganizationID.String(),
			Name:           o.Name,
			Type:           o.Type,
		})
	}
	return items, nil
}

func scopeToLabels(scope string) []string {
	if scope == "full" || scope == "" {
		return []string{"Akses penuh"}
	}
	parts := strings.Split(scope, ",")
	labels := make([]string, 0, len(parts))
	for _, p := range parts {
		switch strings.TrimSpace(p) {
		case "emi":
			labels = append(labels, "EMI")
		case "stability":
			labels = append(labels, "Tren Stabilitas")
		case "risk":
			labels = append(labels, "Risiko")
		default:
			labels = append(labels, p)
		}
	}
	return labels
}

// indexPassports builds a lookup map keyed by passport ID plus the flat list
// of IDs, used to batch-query consents/access logs across all of a user's
// issued passports and enrich each item with its originating passport info.
func indexPassports(passports []entity.IncomePassport) (map[uuid.UUID]entity.IncomePassport, []uuid.UUID) {
	passportMap := make(map[uuid.UUID]entity.IncomePassport, len(passports))
	passportIDs := make([]uuid.UUID, 0, len(passports))
	for _, p := range passports {
		passportMap[p.IncomePassportID] = p
		passportIDs = append(passportIDs, p.IncomePassportID)
	}
	return passportMap, passportIDs
}

func scopeToText(scope []string) string {
	labels := make([]string, 0, len(scope))
	for _, s := range scope {
		switch s {
		case "emi":
			labels = append(labels, "EMI")
		case "stability":
			labels = append(labels, "Tren")
		case "risk":
			labels = append(labels, "Risiko")
		case "full":
			return "Akses penuh"
		default:
			labels = append(labels, s)
		}
	}
	return strings.Join(labels, " + ")
}

func consentStatusLabel(status string, expiresAt *time.Time) string {
	switch status {
	case "revoked":
		return "DICABUT"
	case "expired":
		return "KEDALUWARSA"
	case "active":
		if expiresAt != nil {
			days := int(time.Until(*expiresAt).Hours() / 24)
			if days <= 3 {
				return fmt.Sprintf("%d HR LAGI", days)
			}
		}
		return "AKTIF"
	default:
		return status
	}
}
