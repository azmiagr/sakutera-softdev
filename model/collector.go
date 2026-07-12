package model

type CreatePairingCodeRequest struct {
	ExpiresInSeconds int `json:"expires_in_seconds"`
}

type QRPayload struct {
	Type    string `json:"type"`
	Version int    `json:"version"`
	Code    string `json:"code"`
}

type CreatePairingCodeResponse struct {
	PairingCode string    `json:"pairing_code"`
	ExpiresAt   string    `json:"expires_at"`
	QRPayload   QRPayload `json:"qr_payload"`
}

type PairDeviceRequest struct {
	PairingCode string `json:"pairing_code" binding:"required"`
	DeviceID    string `json:"device_id" binding:"required"`
	DeviceName  string `json:"device_name"`
	Platform    string `json:"platform"`
	OSVersion   string `json:"os_version"`
	AppVersion  string `json:"app_version"`
}

type PairDeviceResponse struct {
	DeviceID    string `json:"device_id"`
	DeviceToken string `json:"device_token"`
	PairedAt    string `json:"paired_at"`
}

type DeviceItem struct {
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	Platform   string `json:"platform"`
	AppVersion string `json:"app_version"`
	IsActive   bool   `json:"is_active"`
	PairedAt   string `json:"paired_at"`
	LastSeenAt string `json:"last_seen_at"`
}

type ListDevicesResponse struct {
	Devices []DeviceItem `json:"devices"`
}

type RevokeDeviceResponse struct {
	DeviceID  string `json:"device_id"`
	IsActive  bool   `json:"is_active"`
	RevokedAt string `json:"revoked_at"`
}

type CollectorHealthResponse struct {
	Connected    bool   `json:"connected"`
	DeviceActive bool   `json:"device_active"`
	ServerTime   string `json:"server_time"`
}

type AllowedPackageItem struct {
	PackageName string `json:"package_name"`
	Provider    string `json:"provider"`
	Enabled     bool   `json:"enabled"`
}

type CollectorConfigResponse struct {
	Mode                string               `json:"mode"`
	ConfigVersion       int                  `json:"config_version"`
	AllowedPackages     []AllowedPackageItem `json:"allowed_packages"`
	MaxBatchSize        int                  `json:"max_batch_size"`
	SyncIntervalMinutes int                  `json:"sync_interval_minutes"`
}

type NotificationEventItem struct {
	EventID        string `json:"event_id" binding:"required"`
	PackageName    string `json:"package_name" binding:"required"`
	NotificationID int64  `json:"notification_id"`
	Title          string `json:"title"`
	Text           string `json:"text"`
	BigText        string `json:"big_text"`
	PostedAt       int64  `json:"posted_at" binding:"required"`
	CapturedAt     int64  `json:"captured_at" binding:"required"`
}

type NotificationBatchRequest struct {
	DeviceID string                  `json:"device_id" binding:"required"`
	Events   []NotificationEventItem `json:"events" binding:"required,min=1,dive"`
}

type NotificationEventResult struct {
	EventID       string  `json:"event_id"`
	Status        string  `json:"status"`
	TransactionID *string `json:"transaction_id"`
	ReviewID      *string `json:"review_id"`
	ErrorCode     *string `json:"error_code"`
}

type NotificationBatchResponse struct {
	Received    int                       `json:"received"`
	Processed   int                       `json:"processed"`
	NeedsReview int                       `json:"needs_review"`
	Duplicates  int                       `json:"duplicates"`
	Ignored     int                       `json:"ignored"`
	Failed      int                       `json:"failed"`
	Results     []NotificationEventResult `json:"results"`
}
