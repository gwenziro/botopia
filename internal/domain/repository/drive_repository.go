package repository

import (
	"context"
	"io"

	"google.golang.org/api/drive/v3"
)

// DriveRepository mendefinisikan kontrak untuk repository penyimpanan file
type DriveRepository interface {
	// UploadFile mengupload file ke penyimpanan
	UploadFile(ctx context.Context, name string, mimeType string, content io.Reader) (string, error)

	// GetFileURL mendapatkan URL file
	GetFileURL(fileID string) string

	// IsConfigured memeriksa apakah repository sudah dikonfigurasi
	IsConfigured() bool

	// UploadImage mengunggah file gambar ke Google Drive
	UploadImage(ctx context.Context, filePath string, transactionCode string) (string, error)

	// GetDriveService mendapatkan akses ke service Google Drive
	GetDriveService(ctx context.Context) (*drive.Service, error)
}
