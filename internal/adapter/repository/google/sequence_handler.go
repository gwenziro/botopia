package google

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gwenziro/botopia/internal/domain/finance"
	"github.com/gwenziro/botopia/internal/infrastructure/config"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
	"google.golang.org/api/sheets/v4"
)

// SequenceHandler menangani operasi untuk penomoran urut
type SequenceHandler struct {
	apiRepo *GoogleAPIRepository
	config  *config.GoogleSheetsConfig
	log     *logger.Logger
}

// NewSequenceHandler membuat instance sequence handler baru
func NewSequenceHandler(
	apiRepo *GoogleAPIRepository,
	config *config.GoogleSheetsConfig,
	log *logger.Logger,
) *SequenceHandler {
	return &SequenceHandler{
		apiRepo: apiRepo,
		config:  config,
		log:     log,
	}
}

// GetNextSequenceNumber mendapatkan nomor urut berikutnya dalam bulan ini (untuk kode unik)
func (h *SequenceHandler) GetNextSequenceNumber(_ context.Context, service *sheets.Service, sheetName string) (int, error) {
	// Ambil data bulan ini
	now := time.Now()
	currentMonth := finance.GetMonthAbbr(now.Month())
	currentYear := now.Year() % 100
	prefix := "k"
	if sheetName == "Pemasukan" {
		prefix = "m"
	}

	// Pattern untuk bulan dan tahun ini: k_mei25_ atau m_mei25_
	pattern := fmt.Sprintf("%s_%s%02d_", prefix, currentMonth, currentYear)

	// Ambil semua data
	resp, err := service.Spreadsheets.Values.Get(
		h.config.SpreadsheetID,
		sheetName+"!A:B", // Kolom A & B (No dan Kode Unik)
	).Do()

	if err != nil {
		return 0, err
	}

	maxSeq := 0

	// Cari nomor urut maksimal untuk bulan ini
	if len(resp.Values) > 0 {
		for _, row := range resp.Values {
			if len(row) >= 2 && row[1] != nil {
				codeStr := fmt.Sprintf("%v", row[1])
				if strings.HasPrefix(codeStr, pattern) {
					// Extract nomor urut dari kode, format: k_mei25_002
					seqPart := codeStr[len(pattern):] // Ambil "002"
					seq, err := strconv.Atoi(seqPart)
					if err == nil && seq > maxSeq {
						maxSeq = seq
					}
				}
			}
		}
	}

	// Kembalikan nomor berikutnya
	return maxSeq + 1, nil
}

// GetGlobalRecordNumber mendapatkan nomor urut global untuk sheet tertentu
func (h *SequenceHandler) GetGlobalRecordNumber(_ context.Context, service *sheets.Service, sheetName string) (int, error) {
	// Ambil semua data
	resp, err := service.Spreadsheets.Values.Get(
		h.config.SpreadsheetID,
		sheetName+"!A:A", // Hanya kolom A (nomor)
	).Do()

	if err != nil {
		return 0, err
	}

	// Hitung jumlah baris data (dikurangi header)
	rowCount := len(resp.Values)
	if rowCount > 0 {
		// Kembalikan nomor berikutnya setelah header
		return rowCount, nil
	}

	// Jika belum ada data, mulai dari 1
	return 1, nil
}
