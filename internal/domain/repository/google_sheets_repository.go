package repository

import (
	"context"

	"github.com/gwenziro/botopia/internal/domain/finance"
	"google.golang.org/api/sheets/v4"
)

// GoogleSheetsRepository mendefinisikan kontrak untuk akses Google Sheets
type GoogleSheetsRepository interface {
	// AddIncomeRecord menambahkan record pemasukan ke sheet
	AddIncomeRecord(ctx context.Context, record *finance.FinanceRecord) error

	// AddExpenseRecord menambahkan record pengeluaran ke sheet
	AddExpenseRecord(ctx context.Context, record *finance.FinanceRecord) error

	// GetRecentRecords mendapatkan record terbaru (gabungan pemasukan & pengeluaran)
	GetRecentRecords(ctx context.Context, limit int) ([]*finance.FinanceRecord, error)

	// GetConfiguration mendapatkan konfigurasi dari sheet
	GetConfiguration(ctx context.Context) (*finance.Configuration, error)

	// GetSpreadsheetURL mendapatkan URL spreadsheet
	GetSpreadsheetURL() string

	// GetSheetsService mendapatkan akses ke service Google Sheets
	GetSheetsService(ctx context.Context) (*sheets.Service, error)

	// IsConfigured memeriksa apakah repository sudah dikonfigurasi
	IsConfigured() bool
}
