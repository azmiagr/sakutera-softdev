package model

type PassportEligibility struct {
	IsEligible      bool  `json:"is_eligible"`
	DaysOfData      int   `json:"days_of_data"`
	MinRequired     int   `json:"min_required"`
	EntriesVerified int64 `json:"entries_verified"`
}

type ActivePassportItem struct {
	IncomePassportID string  `json:"income_passport_id"`
	PassportNumber   string  `json:"passport_number"`
	EMIValue         float64 `json:"emi_value"`
	PeriodType       string  `json:"period_type"`
	PeriodLabel      string  `json:"period_label"`
	RiskLevel        string  `json:"risk_level"`
	IssuedAt         string  `json:"issued_at"`
}

type GetPassportResponse struct {
	Eligibility    PassportEligibility `json:"eligibility"`
	ActivePassport *ActivePassportItem `json:"active_passport"`
}

type PassportPreviewResponse struct {
	PeriodType     string  `json:"period_type"`
	PeriodLabel    string  `json:"period_label"`
	PeriodStart    string  `json:"period_start"`
	PeriodEnd      string  `json:"period_end"`
	EMIValue       float64 `json:"emi_value"`
	StabilityLabel string  `json:"stability_label"`
	TrendDirection string  `json:"trend_direction"`
	TrendChangePct float64 `json:"trend_change_pct"`
	RiskLevel      string  `json:"risk_level"`
	RiskScore      float64 `json:"risk_score"`
	TotalEntries   int     `json:"total_entries"`
}

type IssuePassportRequest struct {
	Period string `json:"period" binding:"required,oneof=3_bulan 6_bulan 12_bulan"`
}

type IssuePassportResponse struct {
	IncomePassportID string  `json:"income_passport_id"`
	PassportNumber   string  `json:"passport_number"`
	EMIValue         float64 `json:"emi_value"`
	PeriodType       string  `json:"period_type"`
	PeriodLabel      string  `json:"period_label"`
	RiskLevel        string  `json:"risk_level"`
	IssuedAt         string  `json:"issued_at"`
}
