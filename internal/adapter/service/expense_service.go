// New file for expense-specific service methods
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gwenziro/botopia/internal/domain/finance"
)

// AddExpense menambahkan record pengeluaran baru
func (s *FinanceService) AddExpense(
	ctx context.Context,
	description string,
	amount float64,
	category string,
	paymentMethod string,
	storageMedia string,
	notes string,
	proofURL string,
) (*finance.FinanceRecord, error) {
	// Use current date
	return s.AddExpenseWithDate(
		ctx,
		time.Now(),
		description,
		amount,
		category,
		paymentMethod,
		storageMedia,
		notes,
		proofURL,
	)
}

// AddExpenseWithDate menambahkan record pengeluaran baru dengan tanggal kustom
func (s *FinanceService) AddExpenseWithDate(
	ctx context.Context,
	date time.Time,
	description string,
	amount float64,
	category string,
	paymentMethod string,
	storageMedia string,
	notes string,
	proofURL string,
) (*finance.FinanceRecord, error) {
	s.log.Info("Menambahkan pengeluaran baru dengan tanggal kustom: %s - %s (%.2f)",
		date.Format("2006-01-02"), description, amount)

	// Handle blank notes
	if notes == "" {
		notes = "-"
	}

	// Validasi kategori, metode pembayaran, dan sumber dana terhadap konfigurasi
	if s.config != nil {
		if !contains(s.config.ExpenseCategories, category) {
			return nil, fmt.Errorf("kategori pengeluaran '%s' tidak valid", category)
		}
		if !contains(s.config.PaymentMethods, paymentMethod) {
			return nil, fmt.Errorf("metode pembayaran '%s' tidak valid", paymentMethod)
		}
		if !contains(s.config.StorageMedias, storageMedia) {
			return nil, fmt.Errorf("sumber dana '%s' tidak valid", storageMedia)
		}
	}

	// Buat record untuk pengeluaran
	record := &finance.FinanceRecord{
		Date:          date,
		Description:   description,
		Amount:        amount,
		Category:      category,
		PaymentMethod: paymentMethod,
		StorageMedia:  storageMedia,
		Notes:         notes,
		ProofURL:      proofURL,
		Type:          finance.TypeExpense,
	}

	// Validasi record
	if err := record.Validate(); err != nil {
		return nil, err
	}

	// Tambahkan ke Google Sheets
	err := s.sheetsRepo.AddExpenseRecord(ctx, record)
	if err != nil {
		s.log.Error("Gagal menambahkan pengeluaran ke sheet: %v", err)
		return nil, fmt.Errorf("gagal menambahkan pengeluaran: %v", err)
	}

	s.log.Info("Pengeluaran berhasil dicatat dengan kode: %s", record.UniqueCode)
	return record, nil
}

// ValidateAddExpenseParams memvalidasi parameter untuk penambahan pengeluaran
func (s *FinanceService) ValidateAddExpenseParams(
	ctx context.Context,
	category, paymentMethod, storageMedia string,
) error {
	// Load config if needed
	var err error
	if s.config == nil {
		s.config, err = s.GetConfiguration(ctx)
		if err != nil {
			return fmt.Errorf("gagal memuat konfigurasi: %v", err)
		}
	}

	// Validasi kategori
	if !contains(s.config.ExpenseCategories, category) {
		return fmt.Errorf("kategori '%s' tidak valid. Kategori yang tersedia: %v",
			category, s.config.ExpenseCategories)
	}

	// Validasi metode pembayaran
	if !contains(s.config.PaymentMethods, paymentMethod) {
		return fmt.Errorf("metode pembayaran '%s' tidak valid. Metode yang tersedia: %v",
			paymentMethod, s.config.PaymentMethods)
	}

	// Validasi sumber dana
	if !contains(s.config.StorageMedias, storageMedia) {
		return fmt.Errorf("sumber dana '%s' tidak valid. Sumber dana yang tersedia: %v",
			storageMedia, s.config.StorageMedias)
	}

	return nil
}
