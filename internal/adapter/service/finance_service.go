package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gwenziro/botopia/internal/domain/finance"
	"github.com/gwenziro/botopia/internal/domain/repository"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
)

// FinanceService implementasi layanan keuangan
type FinanceService struct {
	sheetsRepo     repository.FinanceRepository
	driveRepo      repository.DriveRepository
	config         *finance.Configuration
	configErr      error
	configInitOnce sync.Once
	log            *logger.Logger
}

// NewFinanceService membuat instance layanan keuangan baru
func NewFinanceService(
	sheetsRepo repository.FinanceRepository,
	driveRepo repository.DriveRepository,
	log *logger.Logger,
) *FinanceService {
	s := &FinanceService{
		sheetsRepo: sheetsRepo,
		driveRepo:  driveRepo,
		log:        log,
	}

	// Muat konfigurasi di background
	go s.prefetchConfiguration()

	return s
}

// prefetchConfiguration mengambil konfigurasi di background saat service dibuat
func (s *FinanceService) prefetchConfiguration() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.log.Info("Memuat konfigurasi keuangan awal...")

	// Cek apakah repository terkonfigurasi
	if !s.sheetsRepo.IsConfigured() {
		s.log.Warn("Repository Google Sheets tidak terkonfigurasi dengan benar")
		s.configErr = fmt.Errorf("repository tidak terkonfigurasi dengan benar")
		return
	}

	config, err := s.sheetsRepo.GetConfiguration(ctx)
	if err != nil {
		s.log.Warn("Gagal memuat konfigurasi awal: %v", err)
		s.configErr = err
		return
	}

	s.config = config
	s.log.Info("Konfigurasi keuangan berhasil dimuat:")
	s.log.Info("- %d kategori pemasukan", len(config.IncomeCategories))
	s.log.Info("- %d kategori pengeluaran", len(config.ExpenseCategories))
	s.log.Info("- %d media penyimpanan", len(config.StorageMedias))
	s.log.Info("- %d metode pembayaran", len(config.PaymentMethods))

	// Tampilkan URL Spreadsheet
	s.log.Info("Spreadsheet URL: %s", s.GetSpreadsheetURL())
}

// GetRecentRecords mendapatkan record keuangan terbaru
func (s *FinanceService) GetRecentRecords(ctx context.Context, limit int) ([]*finance.FinanceRecord, error) {
	s.log.Info("Mengambil %d record keuangan terbaru", limit)
	return s.sheetsRepo.GetRecentRecords(ctx, limit)
}

// GetSpreadsheetURL mendapatkan URL spreadsheet
func (s *FinanceService) GetSpreadsheetURL() string {
	return s.sheetsRepo.GetSpreadsheetURL()
}

// GetConfiguration mendapatkan konfigurasi keuangan
func (s *FinanceService) GetConfiguration(ctx context.Context) (*finance.Configuration, error) {
	// Gunakan cache jika tersedia
	if s.config != nil && s.configErr == nil {
		return s.config, nil
	}

	// Ambil ulang jika tidak tersedia atau ada error sebelumnya
	config, err := s.sheetsRepo.GetConfiguration(ctx)
	if err != nil {
		return nil, err
	}

	// Update cache
	s.config = config
	s.configErr = nil

	return config, nil
}
