package service

import (
	"context"
	"time"

	"github.com/gwenziro/botopia/internal/domain/finance"
)

// FinanceService mendefinisikan layanan untuk manajemen keuangan
type FinanceService interface {
	// AddIncome menambahkan record pemasukan baru
	AddIncome(ctx context.Context, description string, amount float64, category, storageMedia, notes string, proofURL string) (*finance.FinanceRecord, error)

	// AddExpense menambahkan record pengeluaran baru
	AddExpense(ctx context.Context, description string, amount float64, category, paymentMethod, storageMedia, notes string, proofURL string) (*finance.FinanceRecord, error)

	// AddIncomeWithDate menambahkan record pemasukan baru dengan tanggal kustom
	AddIncomeWithDate(ctx context.Context, date time.Time, description string, amount float64, category, storageMedia, notes string, proofURL string) (*finance.FinanceRecord, error)

	// AddExpenseWithDate menambahkan record pengeluaran baru dengan tanggal kustom
	AddExpenseWithDate(ctx context.Context, date time.Time, description string, amount float64, category, paymentMethod, storageMedia, notes string, proofURL string) (*finance.FinanceRecord, error)

	// GetConfiguration mendapatkan konfigurasi keuangan
	GetConfiguration(ctx context.Context) (*finance.Configuration, error)

	// ValidateAddIncomeParams memvalidasi parameter untuk penambahan pemasukan
	ValidateAddIncomeParams(ctx context.Context, category, storageMedia string) error

	// ValidateAddExpenseParams memvalidasi parameter untuk penambahan pengeluaran
	ValidateAddExpenseParams(ctx context.Context, category, paymentMethod, storageMedia string) error

	// GetSpreadsheetURL mendapatkan URL spreadsheet
	GetSpreadsheetURL() string

	// GetSpreadsheetID mendapatkan ID spreadsheet
	GetSpreadsheetID() string

	// GetRecentRecords mendapatkan record keuangan terbaru
	GetRecentRecords(ctx context.Context, limit int) ([]*finance.FinanceRecord, error)

	// UploadTransactionProof mengunggah bukti transaksi
	UploadTransactionProof(ctx context.Context, transactionCode string, filePath string) (*finance.FinanceRecord, error)

	// UpdateConfiguration memperbarui konfigurasi keuangan
	UpdateConfiguration(ctx context.Context, config *finance.Configuration) error
}
