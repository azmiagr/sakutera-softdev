package entity

import (
	"time"

	"github.com/google/uuid"
)

type ForecastResult struct {
	ForecastResultID     uuid.UUID  `json:"forecast_result_id" gorm:"type:varchar(36);primaryKey"`
	UserID               uuid.UUID  `json:"user_id" gorm:"type:varchar(36);index"`
	EMIValue             float64    `json:"emi_value" gorm:"type:decimal(15,2)"`
	Confidence           string     `json:"confidence" gorm:"type:enum('low','medium','high');not null"`
	TrendDirection       string     `json:"trend_direction" gorm:"type:enum('up','down','stable');not null"`
	TrendChangePct       float64    `json:"trend_change_pct" gorm:"type:decimal(5,2)"`
	MonthToDateIncome    float64    `json:"month_to_date_income" gorm:"type:decimal(15,2)"`
	ForecastTotal        float64    `json:"forecast_total" gorm:"type:decimal(15,2)"`
	RiskLevel            string     `json:"risk_level" gorm:"type:enum('RENDAH','SEDANG','TINGGI');not null"`
	RiskScore            float64    `json:"risk_score" gorm:"type:decimal(5,4)"`
	ProjectedDeficitDate *time.Time `json:"projected_deficit_date" gorm:"type:date"`
	EstimatedShortfall   float64    `json:"estimated_shortfall" gorm:"type:decimal(15,2)"`
	ShouldNotify         bool       `json:"should_notify" gorm:"type:boolean;default:false;not null"`
	IsDataSufficient     bool       `json:"is_data_sufficient" gorm:"type:boolean;default:false;not null"`
	DaysOfData           int        `json:"days_of_data" gorm:"type:int"`
	TransactionCount     int        `json:"transaction_count" gorm:"type:int"`
	ModelName            string     `json:"model_name" gorm:"type:varchar(30)"`
	ModelVersion         string     `json:"model_version" gorm:"type:varchar(20)"`
	CreatedAt            time.Time  `json:"created_at" gorm:"autoCreateTime"`
}
