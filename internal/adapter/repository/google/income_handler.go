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

// incomeHandler menangani operasi untuk income
type incomeHandler struct {
	apiRepo    *GoogleAPIRepository
	config     *config.GoogleSheetsConfig
	seqHandler *sequenceHandler
	log        *logger.Logger
}

// newIncomeHandler membuat instance income handler baru
func newIncomeHandler(
	apiRepo *GoogleAPIRepository,
	config *config.GoogleSheetsConfig,
	seqHandler *sequenceHandler,
	log *logger.Logger,
) *incomeHandler {
	return &incomeHandler{
		apiRepo:    apiRepo,
		config:     config,
		seqHandler: seqHandler,
		log:        log,
	}
}

// AddRecord menambahkan record pemasukan ke sheet
func (h *incomeHandler) AddRecord(ctx context.Context, record *finance.FinanceRecord) error {
	h.log.Info("Memulai penambahan record pemasukan...")

	service, err := h.apiRepo.GetSheetsService(ctx)
	if err != nil {
		h.log.Error("Gagal mendapatkan sheets service: %v", err)
		return fmt.Errorf("gagal mendapatkan sheets service: %v", err)
	}

	// Validasi record
	if err := record.Validate(); err != nil {
		h.log.Error("Validasi record gagal: %v", err)
		return err
	}
	h.log.Debug("Validasi record berhasil")

	// Dapatkan nomor urut dan buat kode unik
	sequenceNumber, err := h.seqHandler.GetNextSequenceNumber(ctx, service, "Pemasukan")
	if err != nil {
		h.log.Error("Gagal mendapatkan nomor urut: %v", err)
		return fmt.Errorf("gagal mendapatkan nomor urut: %v", err)
	}
	h.log.Debug("Nomor urut berikutnya: %d", sequenceNumber)

	// Buat kode unik sesuai format: m_mei25_001
	record.UniqueCode = finance.GenerateUniqueCode(
		finance.TypeIncome,
		record.Date,
		sequenceNumber,
	)
	h.log.Debug("Kode unik dibuat: %s", record.UniqueCode)

	// Dapatkan nomor global untuk kolom nomor
	globalNumber, err := h.seqHandler.GetGlobalRecordNumber(ctx, service, "Pemasukan")
	if err != nil {
		h.log.Error("Gagal mendapatkan nomor global: %v", err)
		return fmt.Errorf("gagal mendapatkan nomor global: %v", err)
	}

	// Buat row baru
	values := []interface{}{
		globalNumber,                     // No urut global
		record.UniqueCode,                // Kode Unik sesuai format
		record.Date.Format("02/01/2006"), // Format tanggal DD/MM/YYYY
		record.Description,               // Deskripsi
		"",                               // Deskripsi merged dengan kolom sebelumnya
		record.Amount,                    // Nominal
		record.Category,                  // Kategori
		record.StorageMedia,              // Media Penyimpanan
		record.Notes,                     // Keterangan (opsional)
		record.ProofURL,                  // Bukti URL (opsional)
	}

	// Append ke sheet Pemasukan
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{values},
	}

	_, err = service.Spreadsheets.Values.Append(
		h.config.SpreadsheetID,
		"Pemasukan!A:J", // Range sesuai struktur sheet
		valueRange,
	).ValueInputOption("USER_ENTERED").Do()

	if err != nil {
		h.log.Error("Gagal menambahkan data ke sheet: %v", err)
		return err
	}

	h.log.Info("Record pemasukan berhasil ditambahkan dengan kode: %s", record.UniqueCode)
	return nil
}

// GetRecords mendapatkan semua record pemasukan
func (h *incomeHandler) GetRecords(ctx context.Context) ([]*finance.FinanceRecord, error) {
	service, err := h.apiRepo.GetSheetsService(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan sheets service: %v", err)
	}

	var records []*finance.FinanceRecord

	// Ambil data pemasukan
	incomeResp, err := service.Spreadsheets.Values.Get(
		h.config.SpreadsheetID,
		"Pemasukan!A2:J", // Skip header row
	).Do()
	if err != nil {
		return nil, fmt.Errorf("gagal membaca data pemasukan: %v", err)
	}

	// Proses data pemasukan
	if len(incomeResp.Values) > 0 {
		for _, row := range incomeResp.Values {
			if len(row) < 9 {
				continue
			}

			record, err := h.parseRow(row)
			if err != nil {
				h.log.Warn("Gagal parse row pemasukan: %v", err)
				continue
			}

			records = append(records, record)
		}
	}

	return records, nil
}

// parseRow mengkonversi baris sheet menjadi FinanceRecord untuk pemasukan
func (h *incomeHandler) parseRow(row []interface{}) (*finance.FinanceRecord, error) {
	record := &finance.FinanceRecord{
		Type: finance.TypeIncome,
	}

	// Parse the row values
	if len(row) > 0 && row[0] != nil {
		if num, err := strconv.Atoi(fmt.Sprintf("%v", row[0])); err == nil {
			record.Number = num
		}
	}

	if len(row) > 1 && row[1] != nil {
		record.UniqueCode = fmt.Sprintf("%v", row[1])
	}

	if len(row) > 2 && row[2] != nil {
		dateStr := fmt.Sprintf("%v", row[2])
		date, err := time.Parse("02/01/2006", dateStr)
		if err != nil {
			// Try alternate format
			date, err = time.Parse("2006-01-02", dateStr)
			if err != nil {
				return nil, fmt.Errorf("invalid date format: %s", dateStr)
			}
		}
		record.Date = date
	}

	if len(row) > 3 && row[3] != nil {
		record.Description = fmt.Sprintf("%v", row[3])
	}

	if len(row) > 5 && row[5] != nil {
		amount, err := strconv.ParseFloat(strings.Replace(fmt.Sprintf("%v", row[5]), ",", ".", -1), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid amount: %v", row[5])
		}
		record.Amount = amount
	}

	if len(row) > 6 && row[6] != nil {
		record.Category = fmt.Sprintf("%v", row[6])
	}

	if len(row) > 7 && row[7] != nil {
		record.StorageMedia = fmt.Sprintf("%v", row[7])
	}

	if len(row) > 8 && row[8] != nil {
		record.Notes = fmt.Sprintf("%v", row[8])
	}

	if len(row) > 9 && row[9] != nil {
		record.ProofURL = fmt.Sprintf("%v", row[9])
	}

	return record, nil
}
