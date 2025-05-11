package repository

import (
	"context"
	"io"

	"google.golang.org/api/drive/v3"
)

// GoogleDriveRepository mendefinisikan kontrak untuk akses Google Drive
type GoogleDriveRepository interface {
	// UploadFile mengupload file ke Google Drive
	UploadFile(ctx context.Context, name string, mimeType string, content io.Reader) (string, error)

	// GetFileURL mendapatkan URL file
	GetFileURL(fileID string) string

	// GetDriveService mendapatkan akses ke service Google Drive
	GetDriveService(ctx context.Context) (*drive.Service, error)

	// IsConfigured memeriksa apakah repository sudah dikonfigurasi
	IsConfigured() bool
}
