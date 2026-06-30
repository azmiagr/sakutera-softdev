package model

type GrantAccessRequest struct {
	OrganizationID string   `json:"organization_id" binding:"required"`
	DataScope      []string `json:"data_scope" binding:"required,min=1"`
	ExpiresInDays  int      `json:"expires_in_days"`
	Purpose        string   `json:"purpose"`
}

type ConsentItem struct {
	ConsentID        string   `json:"consent_id"`
	OrganizationName string   `json:"organization_name"`
	OrganizationType string   `json:"organization_type"`
	GrantedAt        string   `json:"granted_at"`
	DataScope        []string `json:"data_scope"`
	ExpiresAt        *string  `json:"expires_at"`
	DaysRemaining    *int     `json:"days_remaining"`
	Status           string   `json:"status"`
	StatusLabel      string   `json:"status_label"`
	Purpose          string   `json:"purpose"`
}

type GetConsentsResponse struct {
	Consents []ConsentItem `json:"consents"`
}

type AccessLogItem struct {
	AccessLogID      string   `json:"access_log_id"`
	OrganizationName string   `json:"organization_name"`
	OrganizationType string   `json:"organization_type"`
	AccessedAt       string   `json:"accessed_at"`
	DataScope        []string `json:"data_scope"`
	ConsentStatus    string   `json:"consent_status"`
	StatusLabel      string   `json:"status_label"`
	Note             string   `json:"note"`
}

type GetAccessLogsResponse struct {
	Logs []AccessLogItem `json:"logs"`
}

type OrganizationItem struct {
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	Type           string `json:"type"`
}
