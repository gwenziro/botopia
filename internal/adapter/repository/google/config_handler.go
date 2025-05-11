package google

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/gwenziro/botopia/internal/domain/finance"
	"github.com/gwenziro/botopia/internal/infrastructure/config"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
)

// configHandler menangani operasi untuk konfigurasi
type configHandler struct {
	apiRepo *GoogleAPIRepository
	config  *config.GoogleSheetsConfig
	log     *logger.Logger
}

// newConfigHandler membuat instance config handler baru
func newConfigHandler(
	apiRepo *GoogleAPIRepository,
	config *config.GoogleSheetsConfig,
	log *logger.Logger,
) *configHandler {
	return &configHandler{
		apiRepo: apiRepo,
		config:  config,
		log:     log,
	}
}

// GetConfiguration mendapatkan konfigurasi dari sheet
func (h *configHandler) GetConfiguration(ctx context.Context) (*finance.Configuration, error) {
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
func (h *configHandler) extractConfigValues(values [][]interface{}, config *finance.Configuration) {
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
func (h *configHandler) SortAndLimitRecords(records []*finance.FinanceRecord, limit int) []*finance.FinanceRecord {
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
