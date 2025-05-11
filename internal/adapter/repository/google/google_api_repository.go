package google

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/gwenziro/botopia/internal/infrastructure/config"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Scopes yang dibutuhkan untuk API
var scopes = []string{
	sheets.SpreadsheetsScope,
	drive.DriveFileScope,
}

// GoogleAPIRepository implementasi repository untuk Google API
type GoogleAPIRepository struct {
	config *config.GoogleSheetsConfig
	log    *logger.Logger
}

// NewGoogleAPIRepository membuat instance repository baru
func NewGoogleAPIRepository(config *config.Config, log *logger.Logger) *GoogleAPIRepository {
	return &GoogleAPIRepository{
		config: config.GoogleSheets,
		log:    log,
	}
}

// GetSheetsService mendapatkan service Google Sheets menggunakan Service Account
func (r *GoogleAPIRepository) GetSheetsService(ctx context.Context) (*sheets.Service, error) {
	r.log.Info("Mempersiapkan Google Sheets Service dengan Service Account...")

	// Baca file Service Account credentials
	serviceAccountJSON, err := ioutil.ReadFile(r.config.CredentialsFile)
	if err != nil {
		r.log.Error("Tidak dapat membaca file kredensial: %v", err)
		return nil, fmt.Errorf("tidak dapat membaca file kredensial: %v", err)
	}
	r.log.Debug("Berhasil membaca kredensial dari %s", r.config.CredentialsFile)

	// Buat config dari Service Account
	config, err := google.JWTConfigFromJSON(serviceAccountJSON, scopes...)
	if err != nil {
		r.log.Error("Gagal membuat konfigurasi JWT: %v", err)
		return nil, fmt.Errorf("gagal membuat konfigurasi JWT: %v", err)
	}
	r.log.Debug("JWT konfigurasi berhasil dibuat untuk email: %s", config.Email)

	// Buat client HTTP dari config
	client := config.Client(ctx)
	r.log.Debug("HTTP Client berhasil dibuat")

	// Inisialisasi service Sheets
	sheetsService, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		r.log.Error("Gagal membuat layanan Sheets: %v", err)
		return nil, fmt.Errorf("gagal membuat layanan Sheets: %v", err)
	}

	r.log.Info("Google Sheets Service berhasil diinisialisasi")
	return sheetsService, nil
}

// GetDriveService mendapatkan service Google Drive menggunakan Service Account
func (r *GoogleAPIRepository) GetDriveService(ctx context.Context) (*drive.Service, error) {
	r.log.Info("Mempersiapkan Google Drive Service dengan Service Account...")

	// Baca file Service Account credentials
	serviceAccountJSON, err := ioutil.ReadFile(r.config.CredentialsFile)
	if err != nil {
		r.log.Error("Tidak dapat membaca file kredensial: %v", err)
		return nil, fmt.Errorf("tidak dapat membaca file kredensial: %v", err)
	}
	r.log.Debug("Berhasil membaca kredensial dari %s", r.config.CredentialsFile)

	// Buat config dari Service Account
	config, err := google.JWTConfigFromJSON(serviceAccountJSON, scopes...)
	if err != nil {
		r.log.Error("Gagal membuat konfigurasi JWT: %v", err)
		return nil, fmt.Errorf("gagal membuat konfigurasi JWT: %v", err)
	}
	r.log.Debug("JWT konfigurasi berhasil dibuat untuk email: %s", config.Email)

	// Buat client HTTP dari config
	client := config.Client(ctx)
	r.log.Debug("HTTP Client berhasil dibuat")

	// Inisialisasi service Drive
	driveService, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		r.log.Error("Gagal membuat layanan Drive: %v", err)
		return nil, fmt.Errorf("gagal membuat layanan Drive: %v", err)
	}

	r.log.Info("Google Drive Service berhasil diinisialisasi")
	return driveService, nil
}

// IsConfigured memeriksa apakah kredensial dan SpreadsheetID sudah tersedia
func (r *GoogleAPIRepository) IsConfigured() bool {
	if r.config == nil || r.config.CredentialsFile == "" || r.config.SpreadsheetID == "" {
		r.log.Warn("Google API tidak terkonfigurasi dengan benar: %v", map[string]string{
			"credentialsFile": r.config.CredentialsFile,
			"spreadsheetID":   r.config.SpreadsheetID,
		})
		return false
	}

	// Cek apakah file credentials dapat dibaca
	_, err := ioutil.ReadFile(r.config.CredentialsFile)
	isConfigured := err == nil

	if isConfigured {
		r.log.Info("Google API terkonfigurasi dengan benar")
		r.log.Info("SpreadsheetID: %s", r.config.SpreadsheetID)
		r.log.Info("DriveFolderID: %s", r.config.DriveFolderID)
	} else {
		r.log.Error("Gagal membaca file kredensial: %v", err)
	}

	return isConfigured
}

// GetCredentialsPath mendapatkan path ke file credentials
func (r *GoogleAPIRepository) GetCredentialsPath() string {
	return r.config.CredentialsFile
}

// SaveToken diimplementasi untuk memenuhi interface GoogleAPIRepository
// Dengan Service Account, sebenarnya tidak perlu menyimpan token
func (r *GoogleAPIRepository) SaveToken(token *oauth2.Token) error {
	r.log.Info("SaveToken dipanggil, tapi tidak dibutuhkan dengan Service Account")
	return nil
}

// GetAuthURL diimplementasi untuk memenuhi interface GoogleAPIRepository
// Dengan Service Account, sebenarnya tidak perlu URL otorisasi
func (r *GoogleAPIRepository) GetAuthURL(redirectURL string) (string, error) {
	r.log.Info("GetAuthURL dipanggil, tapi tidak dibutuhkan dengan Service Account")
	return "", nil
}
