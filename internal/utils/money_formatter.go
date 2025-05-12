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

	// Ganti koma dengan titik untuk format float
	numericStr = strings.Replace(numericStr, ",", ".", -1)

	return strconv.ParseFloat(numericStr, 64)
}
