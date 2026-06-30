package entity

import (
	"time"

	"github.com/google/uuid"
)

type IncomePassport struct {
	IncomePassportID uuid.UUID `json:"income_passport_id" gorm:"type:varchar(36);primaryKey"`
	UserID           uuid.UUID `json:"user_id" gorm:"type:varchar(36);index"`
	PassportNumber   string    `json:"passport_number" gorm:"type:varchar(100);not null;unique"`
	PeriodType       string    `json:"period_type" gorm:"type:enum('3_bulan','6_bulan','12_bulan');not null"`
	PeriodStart      time.Time `json:"period_start" gorm:"type:date;not null"`
	PeriodEnd        time.Time `json:"period_end" gorm:"type:date;not null"`
	EMI              float64   `json:"emi" gorm:"type:decimal(15,2)"`
	EMI3Month        float64   `json:"emi_3_month" gorm:"type:decimal(15,2)"`
	EMI6Month        float64   `json:"emi_6_month" gorm:"type:decimal(15,2)"`
	EMI12Month       float64   `json:"emi_12_month" gorm:"type:decimal(15,2)"`
	StabilityScore   float64   `json:"stability_score" gorm:"type:decimal(5,4)"`
	StabilityLabel   string    `json:"stability_label" gorm:"type:enum('STABIL','CUKUP STABIL','FLUKTUATIF');not null"`
	TrendDirection   string    `json:"trend_direction" gorm:"type:enum('up','down','stable');not null"`
	TrendChangePct   float64   `json:"trend_change_pct" gorm:"type:decimal(5,2)"`
	RiskLevel        string    `json:"risk_level" gorm:"type:enum('RENDAH','SEDANG','TINGGI');not null"`
	RiskScore        float64   `json:"risk_score" gorm:"type:decimal(5,4)"`
	TotalEntries     int       `json:"total_entries" gorm:"type:int;not null"`
	PdfURL           string    `json:"pdf_url" gorm:"type:text"`
	DigitalSignature string    `json:"digital_signature" gorm:"type:text"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`

	AccessLogs []AccessLog `gorm:"foreignKey:PassportID;references:IncomePassportID;constraint:onDelete:CASCADE"`
	Consents   []Consent   `gorm:"foreignKey:PassportID;references:IncomePassportID;constraint:onDelete:CASCADE"`
}
