package repository

import (
	"context"

	"github.com/gwenziro/botopia/internal/domain/finance"
)

// FinanceRepository mendefinisikan kontrak untuk repository keuangan
type FinanceRepository interface {
	// AddExpenseRecord menambahkan record pengeluaran ke penyimpanan
	AddExpenseRecord(ctx context.Context, record *finance.FinanceRecord) error

	// AddIncomeRecord menambahkan record pemasukan ke penyimpanan
	AddIncomeRecord(ctx context.Context, record *finance.FinanceRecord) error

	// GetRecentRecords mendapatkan record keuangan terbaru
	GetRecentRecords(ctx context.Context, limit int) ([]*finance.FinanceRecord, error)

	// GetConfiguration mendapatkan konfigurasi keuangan
	GetConfiguration(ctx context.Context) (*finance.Configuration, error)

	// IsConfigured memeriksa apakah repository sudah dikonfigurasi
	IsConfigured() bool

	// GetSpreadsheetURL mendapatkan URL spreadsheet
	GetSpreadsheetURL() string

	// GetSheetsService mendapatkan akses ke service sheets
	GetSheetsService(ctx context.Context) (interface{}, error)
}
