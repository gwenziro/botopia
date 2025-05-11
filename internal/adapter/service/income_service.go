// New file for income-specific service methods
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gwenziro/botopia/internal/domain/finance"
)

// AddIncome menambahkan record pemasukan baru
func (s *FinanceService) AddIncome(
	ctx context.Context,
	description string,
	amount float64,
	category string,
	storageMedia string,
	notes string,
	proofURL string,
) (*finance.FinanceRecord, error) {
	// Use current date
	return s.AddIncomeWithDate(
		ctx,
		time.Now(),
		description,
		amount,
		category,
		storageMedia,
		notes,
		proofURL,
	)
}

// AddIncomeWithDate menambahkan record pemasukan baru dengan tanggal kustom
func (s *FinanceService) AddIncomeWithDate(
	ctx context.Context,
	date time.Time,
	description string,
	amount float64,
	category string,
	storageMedia string,
	notes string,
	proofURL string,
) (*finance.FinanceRecord, error) {
	s.log.Info("Menambahkan pemasukan baru dengan tanggal kustom: %s - %s (%.2f)",
		date.Format("2006-01-02"), description, amount)

	// Handle blank notes
	if notes == "" {
		notes = "-"
	}

	// Validasi kategori dan media penyimpanan terhadap konfigurasi
	if s.config != nil {
		if !contains(s.config.IncomeCategories, category) {
			return nil, fmt.Errorf("kategori pemasukan '%s' tidak valid", category)
		}
		if !contains(s.config.StorageMedias, storageMedia) {
			return nil, fmt.Errorf("media penyimpanan '%s' tidak valid", storageMedia)
		}
	}

	// Buat record untuk pemasukan
	record := &finance.FinanceRecord{
		Date:         date,
		Description:  description,
		Amount:       amount,
		Category:     category,
		StorageMedia: storageMedia,
		Notes:        notes,
		ProofURL:     proofURL,
		Type:         finance.TypeIncome,
	}

	// Validasi record
	if err := record.Validate(); err != nil {
		return nil, err
	}

	// Tambahkan ke Google Sheets
	err := s.sheetsRepo.AddIncomeRecord(ctx, record)
	if err != nil {
		s.log.Error("Gagal menambahkan pemasukan ke sheet: %v", err)
		return nil, fmt.Errorf("gagal menambahkan pemasukan: %v", err)
	}

	s.log.Info("Pemasukan berhasil dicatat dengan kode: %s", record.UniqueCode)
	return record, nil
}

// ValidateAddIncomeParams memvalidasi parameter untuk penambahan pemasukan
func (s *FinanceService) ValidateAddIncomeParams(
	ctx context.Context,
	category, storageMedia string,
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
	if !contains(s.config.IncomeCategories, category) {
		return fmt.Errorf("kategori '%s' tidak valid. Kategori yang tersedia: %v",
			category, s.config.IncomeCategories)
	}

	// Validasi media penyimpanan
	if !contains(s.config.StorageMedias, storageMedia) {
		return fmt.Errorf("media penyimpanan '%s' tidak valid. Media yang tersedia: %v",
			storageMedia, s.config.StorageMedias)
	}

	return nil
}
