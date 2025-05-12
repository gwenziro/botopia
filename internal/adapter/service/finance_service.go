package service

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gwenziro/botopia/internal/domain/finance"
	"github.com/gwenziro/botopia/internal/domain/repository"
	"github.com/gwenziro/botopia/internal/domain/service"
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

// Memastikan FinanceService mengimplementasikan interface service.FinanceService
var _ service.FinanceService = (*FinanceService)(nil)

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

// UploadTransactionProof mengunggah bukti transaksi
func (s *FinanceService) UploadTransactionProof(ctx context.Context, transactionCode string, filePath string) (*finance.FinanceRecord, error) {
	s.log.Info("Mengunggah bukti transaksi untuk kode: %s, file: %s", transactionCode, filePath)

	// Validasi file
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file bukti transaksi tidak ditemukan")
	}

	// Cari record berdasarkan kode
	record, err := s.findRecordByCode(ctx, transactionCode)
	if err != nil {
		return nil, fmt.Errorf("gagal mencari transaksi: %v", err)
	}

	// Validasi record
	if record == nil {
		return nil, fmt.Errorf("transaksi dengan kode %s tidak ditemukan", transactionCode)
	}

	// Cek apakah sudah ada bukti
	if record.ProofURL != "" && record.ProofURL != "-" {
		return nil, fmt.Errorf("transaksi dengan kode %s sudah memiliki bukti", transactionCode)
	}

	// Unggah ke Google Drive
	fileURL, err := s.driveRepo.UploadImage(ctx, filePath, transactionCode)
	if err != nil {
		return nil, fmt.Errorf("gagal mengunggah bukti: %v", err)
	}

	// Update record dengan URL bukti
	updatedRecord, err := s.updateRecordProof(ctx, record, fileURL)
	if err != nil {
		return nil, fmt.Errorf("gagal memperbarui record: %v", err)
	}

	// Hapus file temporary
	os.Remove(filePath)

	return updatedRecord, nil
}

// findRecordByCode mencari record berdasarkan kode unik
func (s *FinanceService) findRecordByCode(ctx context.Context, code string) (*finance.FinanceRecord, error) {
	// Cek tipe record dari kode (k_ untuk pengeluaran, m_ untuk pemasukan)
	isExpense := strings.HasPrefix(code, "k_")
	isIncome := strings.HasPrefix(code, "m_")

	if !isExpense && !isIncome {
		return nil, fmt.Errorf("format kode tidak valid")
	}

	// Cari di sheet yang sesuai
	return s.sheetsRepo.FindRecordByCode(ctx, code)
}

// updateRecordProof memperbarui record dengan URL bukti
func (s *FinanceService) updateRecordProof(ctx context.Context, record *finance.FinanceRecord, proofURL string) (*finance.FinanceRecord, error) {
	// Update record dengan URL bukti
	record.ProofURL = proofURL

	// Update di Google Sheets
	err := s.sheetsRepo.UpdateRecordProof(ctx, record.UniqueCode, proofURL)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// GetLogger mengembalikan logger service
func (s *FinanceService) GetLogger() *logger.Logger {
	return s.log
}
