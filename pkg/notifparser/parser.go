package notifparser

import "time"

const (
	Version             = "notifparser-v1"
	ConfidenceThreshold = 0.8
)

type Event struct {
	PackageName string
	Title       string
	Text        string
	BigText     string
	PostedAt    time.Time
}

type Candidate struct {
	SourceName      string
	Provider        string
	TransactionType string // "income" | "expense"
	Amount          float64
	Description     string
	TransactionDate time.Time
	Confidence      float64
	Reason          string
}

type parseFunc func(Event) (*Candidate, bool)

var registry = map[string]parseFunc{
	"com.gojek.app":          parseGoPay,
	"com.grabtaxi.passenger": parseGrabOVO,
}

// Parse never returns nil. Confidence 0 means the notification is not
// recognized as a transaction (unknown package or no amount matched).
func Parse(e Event) *Candidate {
	fn, ok := registry[e.PackageName]
	if !ok {
		return &Candidate{Confidence: 0, Reason: "NOT_A_TRANSACTION"}
	}

	candidate, matched := fn(e)
	if !matched {
		return &Candidate{Confidence: 0, Reason: "NOT_A_TRANSACTION"}
	}

	return candidate
}
