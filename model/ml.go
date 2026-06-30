package model

type MLTransaction struct {
	Date     string `json:"date"`
	Amount   int    `json:"amount"`
	Source   string `json:"source"`
	Category string `json:"category,omitempty"`
}

type MLForecastRequest struct {
	UserID       string          `json:"user_id"`
	AsOfDate     string          `json:"as_of_date"`
	Transactions []MLTransaction `json:"transactions"`
}

type MLForecastResponse struct {
	DataSufficiency struct {
		IsSufficient bool `json:"is_sufficient"`
		DaysOfData   int  `json:"days_of_data"`
	} `json:"data_sufficiency"`
	EMI struct {
		Value      float64 `json:"value"`
		Confidence string  `json:"confidence"`
	} `json:"emi"`
	Trend struct {
		Direction string  `json:"direction"`
		ChangePct float64 `json:"change_pct"`
	} `json:"trend"`
	MonthToDate struct {
		TotalEarned float64 `json:"total_earned"`
	} `json:"month_to_date"`
	ForecastRemainingMonth struct {
		ProjectedMonthTotal float64 `json:"projected_month_total"`
	} `json:"forecast_remaining_month"`
	DeficitRisk struct {
		Level                string  `json:"level"`
		Score                float64 `json:"score"`
		ProjectedDeficitDate *string `json:"projected_deficit_date"`
		EstimatedShortfall   float64 `json:"estimated_shortfall"`
		ShouldNotify         bool    `json:"should_notify"`
	} `json:"deficit_risk"`
	TransactionCount int `json:"transaction_count"`
}
