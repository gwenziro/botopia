package google

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gwenziro/botopia/internal/domain/finance"
	"github.com/gwenziro/botopia/internal/infrastructure/config"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
	"google.golang.org/api/sheets/v4"
)

// ConfigHandler menangani operasi untuk konfigurasi dan operasi record umum
type ConfigHandler struct {
	apiRepo *GoogleAPIRepository
	config  *config.GoogleSheetsConfig
	log     *logger.Logger
}

// NewConfigHandler membuat instance config handler baru
func NewConfigHandler(
	apiRepo *GoogleAPIRepository,
	config *config.GoogleSheetsConfig,
	log *logger.Logger,
) *ConfigHandler {
	return &ConfigHandler{
		apiRepo: apiRepo,
		config:  config,
		log:     log,
	}
}

// GetConfiguration mendapatkan konfigurasi dari sheet
func (h *ConfigHandler) GetConfiguration(ctx context.Context) (*finance.Configuration, error) {
	h.log.Info("Mengambil konfigurasi dari spreadsheet...")

	service, err := h.apiRepo.GetSheetsService(ctx)
	if err != nil {
		h.log.Error("Gagal mendapatkan sheets service: %v", err)
		return nil, fmt.Errorf("gagal mendapatkan sheets service: %v", err)
	}

	// Ambil data dari sheet konfigurasi
	configResp, err := service.Spreadsheets.Values.Get(
		h.config.SpreadsheetID,
		"Konfigurasi!A2:F", // Skip header row
	).Do()
	if err != nil {
		h.log.Error("Gagal membaca data konfigurasi: %v", err)
		return nil, fmt.Errorf("gagal membaca data konfigurasi: %v", err)
	}

	config := &finance.Configuration{
		Year:              time.Now().Year(),
		Month:             int(time.Now().Month()),
		StorageMedias:     []string{},
		PaymentMethods:    []string{},
		ExpenseCategories: []string{},
		IncomeCategories:  []string{},
	}

	if len(configResp.Values) > 0 {
		// Ekstrak data konfigurasi dari sheet
		h.extractConfigValues(configResp.Values, config)
	}

	h.log.Info("Konfigurasi berhasil diambil: %d kategori pemasukan, %d kategori pengeluaran, %d media penyimpanan, %d metode pembayaran",
		len(config.IncomeCategories),
		len(config.ExpenseCategories),
		len(config.StorageMedias),
		len(config.PaymentMethods))

	return config, nil
}

// extractConfigValues mengambil nilai-nilai konfigurasi dari data sheet
func (h *ConfigHandler) extractConfigValues(values [][]interface{}, config *finance.Configuration) {
	// Parse year and month if available
	if len(values[0]) > 0 && values[0][0] != nil {
		year, err := strconv.Atoi(fmt.Sprintf("%v", values[0][0]))
		if err == nil {
			config.Year = year
		}
	}

	if len(values[0]) > 1 && values[0][1] != nil {
		month, err := strconv.Atoi(fmt.Sprintf("%v", values[0][1]))
		if err == nil {
			config.Month = month
		}
	}

	// Extract unique values for each column
	storageMedias := make(map[string]bool)
	paymentMethods := make(map[string]bool)
	expenseCategories := make(map[string]bool)
	incomeCategories := make(map[string]bool)

	for _, row := range values {
		// Media Penyimpanan (column C)
		if len(row) > 2 && row[2] != nil && fmt.Sprintf("%v", row[2]) != "" {
			storageMedias[fmt.Sprintf("%v", row[2])] = true
		}

		// Metode Pembayaran (column D)
		if len(row) > 3 && row[3] != nil && fmt.Sprintf("%v", row[3]) != "" {
			paymentMethods[fmt.Sprintf("%v", row[3])] = true
		}

		// Kategori Pengeluaran (column E)
		if len(row) > 4 && row[4] != nil && fmt.Sprintf("%v", row[4]) != "" {
			expenseCategories[fmt.Sprintf("%v", row[4])] = true
		}

		// Kategori Pemasukan (column F)
		if len(row) > 5 && row[5] != nil && fmt.Sprintf("%v", row[5]) != "" {
			incomeCategories[fmt.Sprintf("%v", row[5])] = true
		}
	}

	// Convert maps to slices
	for k := range storageMedias {
		config.StorageMedias = append(config.StorageMedias, k)
	}

	for k := range paymentMethods {
		config.PaymentMethods = append(config.PaymentMethods, k)
	}

	for k := range expenseCategories {
		config.ExpenseCategories = append(config.ExpenseCategories, k)
	}

	for k := range incomeCategories {
		config.IncomeCategories = append(config.IncomeCategories, k)
	}

	// Sort untuk ordering konsisten
	sort.Strings(config.StorageMedias)
	sort.Strings(config.PaymentMethods)
	sort.Strings(config.ExpenseCategories)
	sort.Strings(config.IncomeCategories)
}

// SortAndLimitRecords mengurutkan dan membatasi jumlah record
func (h *ConfigHandler) SortAndLimitRecords(records []*finance.FinanceRecord, limit int) []*finance.FinanceRecord {
	// Urutkan berdasarkan tanggal terbaru
	sort.Slice(records, func(i, j int) bool {
		return records[i].Date.After(records[j].Date)
	})

	// Batasi jumlah record
	if len(records) > limit {
		records = records[:limit]
	}

	return records
}

// FindRecordByCode mencari record berdasarkan kode unik
func (h *ConfigHandler) FindRecordByCode(ctx context.Context, code string) (*finance.FinanceRecord, error) {
	// Tentukan sheet berdasarkan awalan kode
	sheetName := "Pengeluaran"
	if strings.HasPrefix(code, "m_") {
		sheetName = "Pemasukan"
	}

	// Cari di sheet yang sesuai
	service, err := h.apiRepo.GetSheetsService(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan sheets service: %v", err)
	}

	// Ambil semua data dari sheet
	resp, err := service.Spreadsheets.Values.Get(
		h.config.SpreadsheetID,
		fmt.Sprintf("%s!A2:K", sheetName),
	).Do()
	if err != nil {
		return nil, fmt.Errorf("gagal membaca data sheet: %v", err)
	}

	// Cari baris dengan kode yang cocok
	if len(resp.Values) > 0 {
		for _, row := range resp.Values {
			if len(row) >= 2 && fmt.Sprintf("%v", row[1]) == code {
				// Parse record berdasarkan jenis sheet
				var record *finance.FinanceRecord
				record = &finance.FinanceRecord{
					Type: finance.TypeExpense,
				}
				if sheetName == "Pemasukan" {
					record.Type = finance.TypeIncome
				}

				// Parse basic fields that exist in both sheets
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
					if err == nil {
						record.Amount = amount
					}
				}
				if len(row) > 6 && row[6] != nil {
					record.Category = fmt.Sprintf("%v", row[6])
				}

				// Sheet specific fields
				if sheetName == "Pengeluaran" {
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
				} else { // Pemasukan
					if len(row) > 7 && row[7] != nil {
						record.StorageMedia = fmt.Sprintf("%v", row[7])
					}
					if len(row) > 8 && row[8] != nil {
						record.Notes = fmt.Sprintf("%v", row[8])
					}
					if len(row) > 9 && row[9] != nil {
						record.ProofURL = fmt.Sprintf("%v", row[9])
					}
				}

				return record, nil
			}
		}
	}

	// Record tidak ditemukan
	return nil, nil
}

// UpdateRecordProof memperbarui URL bukti transaksi
func (h *ConfigHandler) UpdateRecordProof(ctx context.Context, code string, proofURL string) error {
	// Tentukan sheet berdasarkan awalan kode
	sheetName := "Pengeluaran"
	if strings.HasPrefix(code, "m_") {
		sheetName = "Pemasukan"
	}

	// Tentukan kolom untuk bukti transaksi
	proofColumn := "K" // Kolom K untuk pengeluaran
	if sheetName == "Pemasukan" {
		proofColumn = "J" // Kolom J untuk pemasukan
	}

	// Cari baris dengan kode yang sesuai
	service, err := h.apiRepo.GetSheetsService(ctx)
	if err != nil {
		return fmt.Errorf("gagal mendapatkan sheets service: %v", err)
	}

	// Ambil semua data dari sheet
	resp, err := service.Spreadsheets.Values.Get(
		h.config.SpreadsheetID,
		fmt.Sprintf("%s!A2:B", sheetName),
	).Do()
	if err != nil {
		return fmt.Errorf("gagal membaca data sheet: %v", err)
	}

	// Cari baris dengan kode yang cocok
	rowIndex := -1
	if len(resp.Values) > 0 {
		for i, row := range resp.Values {
			if len(row) >= 2 && fmt.Sprintf("%v", row[1]) == code {
				rowIndex = i + 2 // +2 karena kita mulai dari A2
				break
			}
		}
	}

	if rowIndex == -1 {
		return fmt.Errorf("record dengan kode %s tidak ditemukan", code)
	}

	// Update cell dengan URL bukti
	updateRange := fmt.Sprintf("%s!%s%d", sheetName, proofColumn, rowIndex)
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{{proofURL}},
	}

	_, err = service.Spreadsheets.Values.Update(
		h.config.SpreadsheetID,
		updateRange,
		valueRange,
	).ValueInputOption("USER_ENTERED").Do()

	if err != nil {
		return fmt.Errorf("gagal memperbarui sheet: %v", err)
	}

	return nil
}
