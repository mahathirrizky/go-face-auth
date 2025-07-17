package helper

import (

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// FormatCurrency formats a float64 amount into a currency string (e.g., IDR).
func FormatCurrency(amount float64) string {
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("Rp%.0f", amount)
}
