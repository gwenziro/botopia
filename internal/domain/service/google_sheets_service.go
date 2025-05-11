package service

import "time"

// FinanceRecord merepresentasikan record keuangan
type FinanceRecord struct {
	ID          string
	Date        time.Time
	Amount      float64
	Category    string
	Description string
	IsIncome    bool
	ImageURL    string // URL Google Drive jika ada bukti
}

// GoogleSheetsService interface untuk interaksi dengan Google Sheets
type GoogleSheetsService interface {
	// AddExpense menambahkan data pengeluaran ke spreadsheet
	AddExpense(amount float64, category, description string, imageURL string) (*FinanceRecord, error)

	// AddIncome menambahkan data pemasukan ke spreadsheet
	AddIncome(amount float64, category, description string, imageURL string) (*FinanceRecord, error)

	// GetLatestRecords mendapatkan 5 record keuangan terakhir
	GetLatestRecords(limit int) ([]*FinanceRecord, error)

	// GetBalance mendapatkan saldo terkini
	GetBalance() (float64, error)

	// GetSheetURL mendapatkan URL spreadsheet
	GetSheetURL() string
}
