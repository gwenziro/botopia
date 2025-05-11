package finance

import (
	"fmt"
	"time"
)

// RecordType definisi tipe record keuangan
type RecordType string

const (
	// TypeIncome adalah record pemasukan
	TypeIncome RecordType = "income"

	// TypeExpense adalah record pengeluaran
	TypeExpense RecordType = "expense"
)

// FinanceRecord merepresentasikan transaksi keuangan
type FinanceRecord struct {
	// Metadata
	Number     int
	UniqueCode string
	Type       RecordType

	// Data Utama
	Date        time.Time
	Description string
	Amount      float64
	Category    string

	// Field khusus pengeluaran
	PaymentMethod string

	// Field umum
	StorageMedia string
	Notes        string
	ProofURL     string
}

// Validate memvalidasi finance record
func (r *FinanceRecord) Validate() error {
	if r.Description == "" {
		return fmt.Errorf("deskripsi harus diisi")
	}

	if r.Amount <= 0 {
		return fmt.Errorf("nominal harus lebih dari 0")
	}

	if r.Category == "" {
		return fmt.Errorf("kategori harus diisi")
	}

	if r.Type == TypeExpense && r.PaymentMethod == "" {
		return fmt.Errorf("metode pembayaran harus diisi untuk pengeluaran")
	}

	if r.StorageMedia == "" {
		return fmt.Errorf("media penyimpanan/sumber dana harus diisi")
	}

	return nil
}

// GetMonthAbbr mengembalikan singkatan bulan dalam bahasa Indonesia
func GetMonthAbbr(month time.Month) string {
	months := []string{
		"jan", "feb", "mar", "apr", "mei", "jun",
		"jul", "agu", "sep", "okt", "nov", "des",
	}

	if int(month) < 1 || int(month) > 12 {
		return "unk"
	}

	return months[month-1]
}

// GenerateUniqueCode membuat kode unik untuk transaksi
func GenerateUniqueCode(typ RecordType, date time.Time, seqNum int) string {
	prefix := "k" // pengeluaran (k from "keluar")
	if typ == TypeIncome {
		prefix = "m" // pemasukan (m from "masuk")
	}

	// Format: x_mmm00_000 (contoh: k_mei23_001)
	monthAbbr := GetMonthAbbr(date.Month())
	yearShort := date.Year() % 100

	return fmt.Sprintf("%s_%s%02d_%03d", prefix, monthAbbr, yearShort, seqNum)
}
