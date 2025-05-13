package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// FormatMoney memformat angka ke format uang dengan pemisah ribuan
func FormatMoney(amount float64) string {
	str := strconv.FormatFloat(amount, 'f', 0, 64)
	result := ""

	for i, c := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += "."
		}
		result += string(c)
	}

	return result
}

// ParseMoney mengkonversi string nominal uang menjadi float64
func ParseMoney(amountStr string) (float64, error) {
	// Bersihkan string dari karakter non-numerik kecuali titik dan koma
	numericStr := regexp.MustCompile(`[^0-9.,]`).ReplaceAllString(amountStr, "")

	// Handle both international and Indonesian formats
	if strings.Contains(numericStr, ",") && strings.Contains(numericStr, ".") {
		// Check which is the decimal separator based on position
		lastCommaPos := strings.LastIndex(numericStr, ",")
		lastDotPos := strings.LastIndex(numericStr, ".")

		if lastCommaPos > lastDotPos {
			// Format like "1.234,56" - comma is decimal separator
			// Remove dots first (thousand separators)
			numericStr = strings.Replace(numericStr, ".", "", -1)
			// Then replace comma with dot for decimal
			numericStr = strings.Replace(numericStr, ",", ".", -1)
		} else {
			// Format like "1,234.56" - dot is decimal separator
			// Just remove commas
			numericStr = strings.Replace(numericStr, ",", "", -1)
		}
	} else if strings.Contains(numericStr, ",") {
		// Only commas, replace with dot
		numericStr = strings.Replace(numericStr, ",", ".", -1)
	}

	return strconv.ParseFloat(numericStr, 64)
}
