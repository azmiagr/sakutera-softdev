package model

type PeriodEligibility struct {
	Period        string `json:"period"`
	Label         string `json:"label"`
	IsEligible    bool   `json:"is_eligible"`
	DaysRequired  int    `json:"days_required"`
	DaysAvailable int    `json:"days_available"`
	RemainingDays int    `json:"remaining_days"`
}

type PassportEligibility struct {
	IsEligible      bool                `json:"is_eligible"`
	DaysOfData      int                 `json:"days_of_data"`
	EntriesVerified int64               `json:"entries_verified"`
	Periods         []PeriodEligibility `json:"periods"`
}

type EligibilityGapDetail struct {
	Period        string `json:"period,omitempty"`
	DaysAvailable int    `json:"days_available"`
	DaysRequired  int    `json:"days_required"`
	RemainingDays int    `json:"remaining_days"`
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
	DaysUsed       int     `json:"days_used"`
	DaysRequired   int     `json:"days_required"`
	IsFullPeriod   bool    `json:"is_full_period"`
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
