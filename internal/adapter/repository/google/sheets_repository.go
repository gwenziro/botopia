package google

import (
	"context"
	"fmt"

	"github.com/gwenziro/botopia/internal/domain/finance"
	"github.com/gwenziro/botopia/internal/infrastructure/config"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
)

// SheetsRepository implementasi Google Sheets repository
type SheetsRepository struct {
	apiRepo        *GoogleAPIRepository
	config         *config.GoogleSheetsConfig
	log            *logger.Logger
	expenseHandler *expenseHandler
	incomeHandler  *incomeHandler
	configHandler  *configHandler
	seqHandler     *sequenceHandler
}

// NewSheetsRepository membuat instance repository baru
func NewSheetsRepository(
	apiRepo *GoogleAPIRepository,
	config *config.Config,
	log *logger.Logger,
) *SheetsRepository {
	repo := &SheetsRepository{
		apiRepo: apiRepo,
		config:  config.GoogleSheets,
		log:     log,
	}

	// Inisialisasi internal handlers
	repo.seqHandler = newSequenceHandler(apiRepo, config.GoogleSheets, log)
	repo.expenseHandler = newExpenseHandler(apiRepo, config.GoogleSheets, repo.seqHandler, log)
	repo.incomeHandler = newIncomeHandler(apiRepo, config.GoogleSheets, repo.seqHandler, log)
	repo.configHandler = newConfigHandler(apiRepo, config.GoogleSheets, log)

	return repo
}

// GetSheetsService mendapatkan akses ke service Google Sheets
// Implementasi ini mengembalikan interface{} untuk sesuai dengan interface FinanceRepository
func (r *SheetsRepository) GetSheetsService(ctx context.Context) (interface{}, error) {
	return r.apiRepo.GetSheetsService(ctx)
}

// IsConfigured memeriksa apakah repository sudah dikonfigurasi
func (r *SheetsRepository) IsConfigured() bool {
	return r.apiRepo.IsConfigured() && r.config.SpreadsheetID != ""
}

// GetSpreadsheetURL mendapatkan URL spreadsheet
func (r *SheetsRepository) GetSpreadsheetURL() string {
	return "https://docs.google.com/spreadsheets/d/" + r.config.SpreadsheetID
}

// AddExpenseRecord menambahkan record pengeluaran ke sheet
func (r *SheetsRepository) AddExpenseRecord(ctx context.Context, record *finance.FinanceRecord) error {
	return r.expenseHandler.AddRecord(ctx, record)
}

// AddIncomeRecord menambahkan record pemasukan ke sheet
func (r *SheetsRepository) AddIncomeRecord(ctx context.Context, record *finance.FinanceRecord) error {
	return r.incomeHandler.AddRecord(ctx, record)
}

// GetConfiguration mendapatkan konfigurasi dari sheet
func (r *SheetsRepository) GetConfiguration(ctx context.Context) (*finance.Configuration, error) {
	return r.configHandler.GetConfiguration(ctx)
}

// GetRecentRecords mendapatkan record terbaru (gabungan pemasukan & pengeluaran)
func (r *SheetsRepository) GetRecentRecords(ctx context.Context, limit int) ([]*finance.FinanceRecord, error) {
	// Ambil data dari kedua handler
	incomeRecords, err := r.incomeHandler.GetRecords(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data pemasukan: %v", err)
	}

	expenseRecords, err := r.expenseHandler.GetRecords(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data pengeluaran: %v", err)
	}

	// Gabung dan urutkan
	records := append(incomeRecords, expenseRecords...)
	return r.configHandler.SortAndLimitRecords(records, limit), nil
}
