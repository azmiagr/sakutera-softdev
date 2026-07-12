package model

type ReviewItem struct {
	ReviewID            string  `json:"review_id"`
	Provider            string  `json:"provider"`
	TransactionType     string  `json:"transaction_type"`
	Amount              float64 `json:"amount"`
	Description         string  `json:"description"`
	TransactionDate     string  `json:"transaction_date"`
	TransactionSourceID *string `json:"transaction_source_id"`
	Confidence          float64 `json:"confidence"`
	Reason              string  `json:"reason"`
	CreatedAt           string  `json:"created_at"`
}

type GetReviewsResponse struct {
	Items      []ReviewItem `json:"items"`
	NextCursor *string      `json:"next_cursor"`
}

type ConfirmReviewRequest struct {
	Amount              float64 `json:"amount"`
	Description         string  `json:"description"`
	TransactionDate     string  `json:"transaction_date"`
	TransactionSourceID string  `json:"transaction_source_id"`
}

type ConfirmReviewResponse struct {
	ReviewID      string `json:"review_id"`
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

type RejectReviewRequest struct {
	Reason string `json:"reason"`
}

type RejectReviewResponse struct {
	ReviewID string `json:"review_id"`
	Status   string `json:"status"`
}
