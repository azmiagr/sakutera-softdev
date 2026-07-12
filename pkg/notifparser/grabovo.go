package notifparser

// parseGrabOVO handles notifications from the Grab app (package placeholder:
// "com.grabtaxi.passenger" — needs verifying against a real device before production use).
func parseGrabOVO(e Event) (*Candidate, bool) {
	text := e.Text
	if e.BigText != "" {
		text = e.BigText
	}
	combined := e.Title + " " + text

	amount, ok := extractAmount(combined)
	if !ok {
		return nil, false
	}

	direction := classifyDirection(combined)

	candidate := &Candidate{
		SourceName:      "Grab",
		Provider:        "OVO",
		Amount:          amount,
		Description:     e.Title,
		TransactionDate: e.PostedAt,
	}

	switch direction {
	case "income":
		candidate.TransactionType = "income"
		candidate.Confidence = 0.9
		candidate.Reason = "amount dan arah income terdeteksi jelas"
	case "expense":
		candidate.TransactionType = "expense"
		candidate.Confidence = 0.9
		candidate.Reason = "notifikasi pengeluaran"
	default:
		candidate.TransactionType = "income"
		candidate.Confidence = 0.4
		candidate.Reason = "amount terdeteksi tapi arah transaksi ambigu"
	}

	return candidate, true
}
