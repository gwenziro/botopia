package google

import (
	"context"
	"fmt"
	"io"

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
