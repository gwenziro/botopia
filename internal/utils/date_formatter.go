package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// IndoMonths adalah nama-nama bulan dalam bahasa Indonesia
var IndoMonths = []string{
	"Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

// IndoMonthsMap adalah pemetaan nama bulan Indonesia ke Inggris
var IndoMonthsMap = map[string]string{
	"jan": "Jan", "feb": "Feb", "mar": "Mar", "apr": "Apr",
	"mei": "May", "jun": "Jun", "jul": "Jul", "agu": "Aug",
	"sep": "Sep", "okt": "Oct", "nov": "Nov", "des": "Dec",
	"januari": "January", "februari": "February", "maret": "March", "april": "April",
	"juni": "June", "juli": "July", "agustus": "August",
	"september": "September", "oktober": "October", "november": "November", "desember": "December",
}

// FormatDateID memformat tanggal ke format Indonesia (DD Bulan YYYY)
func FormatDateID(date time.Time) string {
	return fmt.Sprintf("%02d %s %d", date.Day(), IndoMonths[date.Month()-1], date.Year())
}

// FormatDateShort memformat tanggal ke format pendek DD/MM/YYYY
func FormatDateShort(date time.Time) string {
	return fmt.Sprintf("%02d/%02d/%d", date.Day(), date.Month(), date.Year())
}

// ParseDateWithFormats mencoba mem-parse tanggal dengan multiple format
func ParseDateWithFormats(dateStr string) (time.Time, error) {
	formats := []string{
		"2 Jan 2006",
		"2 January 2006",
		"02 Jan 2006",
		"02 January 2006",
		"2006-01-02",
		"02/01/2006",
		"2/1/2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	// Coba parse format Indonesia (misal: "15 Mei 2025")
	re := regexp.MustCompile(`(\d{1,2})\s+([A-Za-z]+)\s+(\d{4})`)
	match := re.FindStringSubmatch(dateStr)
	if len(match) == 4 {
		day, month, year := match[1], strings.ToLower(match[2]), match[3]
		if englishMonth, ok := IndoMonthsMap[month]; ok {
			newDateStr := fmt.Sprintf("%s %s %s", day, englishMonth, year)
			for _, format := range []string{"2 Jan 2006", "2 January 2006"} {
				if t, err := time.Parse(format, newDateStr); err == nil {
					return t, nil
				}
			}
		}
	}

	// Bila tanggal kosong atau tidak valid, gunakan tanggal hari ini
	if dateStr == "" || dateStr == "hari ini" {
		return time.Now(), nil
	}

	return time.Time{}, fmt.Errorf("format tanggal tidak dikenali")
}
