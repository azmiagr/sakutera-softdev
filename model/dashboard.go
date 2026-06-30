package model

type ForecastData struct {
	EMIValue          float64 `json:"emi_value"`
	Confidence        string  `json:"confidence"`
	TrendDirection    string  `json:"trend_direction"`
	TrendChangePct    float64 `json:"trend_change_pct"`
	MonthToDateIncome float64 `json:"month_to_date_income"`
	ForecastTotal     float64 `json:"forecast_total"`
	RiskLevel         string  `json:"risk_level"`
	RiskScore         float64 `json:"risk_score"`
	IsDataSufficient  bool    `json:"is_data_sufficient"`
	DaysOfData        int     `json:"days_of_data"`
	TransactionCount  int     `json:"transaction_count"`
	ModelName         string  `json:"model_name"`
}

type TrendPoint struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}

type DashboardResponse struct {
	UserFullName       string            `json:"user_full_name"`
	Forecast           *ForecastData     `json:"forecast"`
	TodayIncome        float64           `json:"today_income"`
	ValidChainCount    int64             `json:"valid_chain_count"`
	RecentTransactions []TransactionItem `json:"recent_transactions"`
	TrendData          []TrendPoint      `json:"trend_data"`
}
