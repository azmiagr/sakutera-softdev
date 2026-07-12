package model

type CreateTransactionRequest struct {
	Amount              float64 `json:"amount" binding:"required,gt=0"`
	TransactionSourceID string  `json:"transaction_source_id" binding:"required"`
	TransactionDate     string  `json:"transaction_date" binding:"required"`
	Description         string  `json:"description"`
	AttachmentURL       string  `json:"attachment_url"`
}

type UploadAttachmentRequest struct {
	Amount              float64 `form:"amount"`
	TransactionSourceID string  `form:"transaction_source_id"`
	TransactionDate     string  `form:"transaction_date"`
	Description         string  `form:"description"`
}

type UploadAttachmentResponse struct {
	AttachmentURL       string  `json:"attachment_url"`
	FileType            string  `json:"file_type"`
	Amount              float64 `json:"amount"`
	TransactionSourceID string  `json:"transaction_source_id"`
	SourceName          string  `json:"source_name"`
	SourceProvider      string  `json:"source_provider"`
	TransactionDate     string  `json:"transaction_date"`
	Category            string  `json:"category"`
	Status              string  `json:"status"`
}

type TransactionItem struct {
	TransactionID   string  `json:"transaction_id"`
	SourceName      string  `json:"source_name"`
	SourceProvider  string  `json:"source_provider"`
	Amount          float64 `json:"amount"`
	TransactionDate string  `json:"transaction_date"`
	Category        string  `json:"category"`
	HashPrefix      string  `json:"hash_prefix"`
}

type CreateTransactionResponse struct {
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	CurrentHash   string  `json:"current_hash"`
	Status        string  `json:"status"`
	Message       string  `json:"message"`
}

type TransactionSourceItem struct {
	TransactionSourceID string `json:"transaction_source_id"`
	Name                string `json:"name"`
	Provider            string `json:"provider"`
}

type GetTransactionSourcesResponse struct {
	Sources []TransactionSourceItem `json:"sources"`
}

type GetTransactionsResponse struct {
	Transactions []TransactionItem `json:"transactions"`
	Total        int               `json:"total"`
}
