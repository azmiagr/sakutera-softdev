package notifparser

import (
	"regexp"
	"strconv"
	"strings"
)

var amountRegex = regexp.MustCompile(`Rp\s?([\d.,]+)`)

var incomeKeywords = []string{"menerima", "diterima", "masuk"}
var expenseKeywords = []string{"membayar", "terkirim", "keluar", "bayar", "pembayaran"}

// extractAmount finds the first "Rp<number>" occurrence in the text and
// parses it into a float64, stripping thousands separators.
func extractAmount(text string) (float64, bool) {
	match := amountRegex.FindStringSubmatch(text)
	if match == nil {
		return 0, false
	}

	raw := strings.ReplaceAll(match[1], ".", "")
	raw = strings.ReplaceAll(raw, ",", "")
	raw = strings.TrimRight(raw, ".")

	amount, err := strconv.ParseFloat(raw, 64)
	if err != nil || amount <= 0 {
		return 0, false
	}

	return amount, true
}

// classifyDirection returns "income", "expense", or "" (ambiguous) based on
// keyword matches in the notification text.
func classifyDirection(text string) string {
	lower := strings.ToLower(text)

	isIncome := containsAny(lower, incomeKeywords)
	isExpense := containsAny(lower, expenseKeywords)

	switch {
	case isIncome && !isExpense:
		return "income"
	case isExpense && !isIncome:
		return "expense"
	default:
		return ""
	}
}

func containsAny(text string, keywords []string) bool {
	for _, k := range keywords {
		if strings.Contains(text, k) {
			return true
		}
	}
	return false
}
