package google

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gwenziro/botopia/internal/infrastructure/config"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"

	"google.golang.org/api/drive/v3"
)

// DriveRepository implementasi Google Drive repository
type DriveRepository struct {
	apiRepo *GoogleAPIRepository
	config  *config.GoogleSheetsConfig
	log     *logger.Logger
}

// NewDriveRepository membuat instance repository baru
func NewDriveRepository(
	apiRepo *GoogleAPIRepository,
	config *config.Config,
	log *logger.Logger,
) *DriveRepository {
	return &DriveRepository{
		apiRepo: apiRepo,
		config:  config.GoogleSheets,
		log:     log,
	}
}

// UploadFile mengupload file ke Google Drive
func (r *DriveRepository) UploadFile(
	ctx context.Context,
	name string,
	mimeType string,
	content io.Reader,
) (string, error) {
	service, err := r.apiRepo.GetDriveService(ctx)
	if err != nil {
		return "", fmt.Errorf("gagal mendapatkan drive service: %v", err)
	}

	// Buat file metadata
	file := &drive.File{
		Name:     name,
		MimeType: mimeType,
	}

	// Jika folder ID tersedia, gunakan sebagai parent
	if r.config.DriveFolderID != "" {
		file.Parents = []string{r.config.DriveFolderID}
	}

	// Upload file
	res, err := service.Files.Create(file).
		Media(content).
		Fields("id, webViewLink").
		Do()

	if err != nil {
		return "", fmt.Errorf("gagal upload file: %v", err)
	}

	r.log.Info("File %s berhasil diupload ke Google Drive dengan ID: %s", name, res.Id)

	return res.Id, nil
}

// GetFileURL mendapatkan URL file
func (r *DriveRepository) GetFileURL(fileID string) string {
	return fmt.Sprintf("https://drive.google.com/file/d/%s/view", fileID)
}

// IsConfigured memeriksa apakah repository sudah dikonfigurasi
func (r *DriveRepository) IsConfigured() bool {
	return r.apiRepo.IsConfigured()
}

// UploadImage mengunggah file gambar ke Google Drive
func (r *DriveRepository) UploadImage(ctx context.Context, filePath string, transactionCode string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("gagal membuka file: %v", err)
	}
	defer file.Close()

	// Deteksi jenis mime dari ekstensi file
	mimeType := "image/jpeg" // Default
	fileExt := strings.ToLower(filepath.Ext(filePath))
	switch fileExt {
	case ".png":
		mimeType = "image/png"
	case ".pdf":
		mimeType = "application/pdf"
	case ".mp4":
		mimeType = "video/mp4"
	}

	// Buat nama file yang sesuai
	fileName := fmt.Sprintf("Bukti_%s_%s", transactionCode, filepath.Base(filePath))

	// Upload file dan dapatkan ID
	fileID, err := r.UploadFile(ctx, fileName, mimeType, file)
	if err != nil {
		return "", fmt.Errorf("gagal mengupload file: %v", err)
	}

	// Kembalikan URL file
	return r.GetFileURL(fileID), nil
}

// GetDriveService mendapatkan akses ke service Google Drive
func (r *DriveRepository) GetDriveService(ctx context.Context) (*drive.Service, error) {
	return r.apiRepo.GetDriveService(ctx)
}
