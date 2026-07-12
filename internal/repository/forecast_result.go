package repository

import (
	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IForecastResultRepository interface {
	UpsertByUserID(tx *gorm.DB, result *entity.ForecastResult) error
	GetByUserID(tx *gorm.DB, userID uuid.UUID) (*entity.ForecastResult, error)
}

type ForecastResultRepository struct {
	db *gorm.DB
}

func NewForecastResultRepository(db *gorm.DB) IForecastResultRepository {
	return &ForecastResultRepository{db: db}
}

func (r *ForecastResultRepository) UpsertByUserID(tx *gorm.DB, result *entity.ForecastResult) error {
	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"emi_value", "confidence", "trend_direction", "trend_change_pct",
			"month_to_date_income", "forecast_total", "risk_level", "risk_score",
			"projected_deficit_date", "estimated_shortfall", "should_notify",
			"is_data_sufficient", "days_of_data", "transaction_count",
			"model_name", "model_version", "created_at",
		}),
	}).Create(result).Error
}

func (r *ForecastResultRepository) GetByUserID(tx *gorm.DB, userID uuid.UUID) (*entity.ForecastResult, error) {
	var result entity.ForecastResult
	err := tx.Where("user_id = ?", userID).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
