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

// ExpenseHandler menangani operasi untuk expense
type ExpenseHandler struct {
	apiRepo    *GoogleAPIRepository
	config     *config.GoogleSheetsConfig
	seqHandler *SequenceHandler
	log        *logger.Logger
}

// NewExpenseHandler membuat instance expense handler baru
func NewExpenseHandler(
	apiRepo *GoogleAPIRepository,
	config *config.GoogleSheetsConfig,
	seqHandler *SequenceHandler,
	log *logger.Logger,
) *ExpenseHandler {
	return &ExpenseHandler{
		apiRepo:    apiRepo,
		config:     config,
		seqHandler: seqHandler,
		log:        log,
	}
}

// AddRecord menambahkan record pengeluaran ke sheet
func (h *ExpenseHandler) AddRecord(ctx context.Context, record *finance.FinanceRecord) error {
	h.log.Info("Memulai penambahan record pengeluaran...")

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

	// Dapatkan nomor urut untuk kode unik
	sequenceNumber, err := h.seqHandler.GetNextSequenceNumber(ctx, service, "Pengeluaran")
	if err != nil {
		h.log.Error("Gagal mendapatkan nomor urut: %v", err)
		return fmt.Errorf("gagal mendapatkan nomor urut: %v", err)
	}
	h.log.Debug("Nomor urut berikutnya: %d", sequenceNumber)

	// Buat kode unik sesuai format: k_mei25_002
	record.UniqueCode = finance.GenerateUniqueCode(
		finance.TypeExpense,
		record.Date,
		sequenceNumber,
	)
	h.log.Debug("Kode unik dibuat: %s", record.UniqueCode)

	// Dapatkan nomor global untuk kolom nomor
	globalNumber, err := h.seqHandler.GetGlobalRecordNumber(ctx, service, "Pengeluaran")
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
		record.PaymentMethod,             // Metode Pembayaran
		record.StorageMedia,              // Sumber Dana
		record.Notes,                     // Keterangan (opsional)
		record.ProofURL,                  // Bukti URL (opsional)
	}

	// Append ke sheet Pengeluaran
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{values},
	}

	_, err = service.Spreadsheets.Values.Append(
		h.config.SpreadsheetID,
		"Pengeluaran!A:K", // Range sesuai struktur sheet
		valueRange,
	).ValueInputOption("USER_ENTERED").Do()

	if err != nil {
		h.log.Error("Gagal menambahkan data ke sheet: %v", err)
		return err
	}

	h.log.Info("Record pengeluaran berhasil ditambahkan dengan kode: %s", record.UniqueCode)
	return nil
}

// GetRecords mendapatkan semua record pengeluaran
func (h *ExpenseHandler) GetRecords(ctx context.Context) ([]*finance.FinanceRecord, error) {
	service, err := h.apiRepo.GetSheetsService(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan sheets service: %v", err)
	}

	var records []*finance.FinanceRecord

	// Ambil data pengeluaran
	expenseResp, err := service.Spreadsheets.Values.Get(
		h.config.SpreadsheetID,
		"Pengeluaran!A2:K", // Skip header row
	).Do()
	if err != nil {
		return nil, fmt.Errorf("gagal membaca data pengeluaran: %v", err)
	}

	// Proses data pengeluaran
	if len(expenseResp.Values) > 0 {
		for _, row := range expenseResp.Values {
			if len(row) < 10 {
				continue
			}

			record, err := h.parseRow(row)
			if err != nil {
				h.log.Warn("Gagal parse row pengeluaran: %v", err)
				continue
			}

			records = append(records, record)
		}
	}

	return records, nil
}

// parseRow mengkonversi baris sheet menjadi FinanceRecord untuk pengeluaran
func (h *ExpenseHandler) parseRow(row []interface{}) (*finance.FinanceRecord, error) {
	record := &finance.FinanceRecord{
		Type: finance.TypeExpense,
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
		record.PaymentMethod = fmt.Sprintf("%v", row[7])
	}

	if len(row) > 8 && row[8] != nil {
		record.StorageMedia = fmt.Sprintf("%v", row[8])
	}

	if len(row) > 9 && row[9] != nil {
		record.Notes = fmt.Sprintf("%v", row[9])
	}

	if len(row) > 10 && row[10] != nil {
		record.ProofURL = fmt.Sprintf("%v", row[10])
	}

	return record, nil
}
