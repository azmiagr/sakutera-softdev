package model

type LedgerSummary struct {
	TotalEntries int64 `json:"total_entries"`
	ChainValid   bool  `json:"chain_valid"`
}

type LedgerTransactionItem struct {
	TransactionID   string  `json:"transaction_id"`
	SourceName      string  `json:"source_name"`
	SourceProvider  string  `json:"source_provider"`
	Amount          float64 `json:"amount"`
	TransactionDate string  `json:"transaction_date"`
	TransactionTime string  `json:"transaction_time"`
	IsVerified      bool    `json:"is_verified"`
	HashPrefix      string  `json:"hash_prefix"`
}

type GetLedgerResponse struct {
	Summary      LedgerSummary          `json:"summary"`
	Transactions []LedgerTransactionItem `json:"transactions"`
	Total        int                    `json:"total"`
}
