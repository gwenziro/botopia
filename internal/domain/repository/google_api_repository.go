package repository

import (
	"context"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
)

// GoogleAPIRepository mendefinisikan kontrak untuk akses Google API
type GoogleAPIRepository interface {
	// GetSheetsService mendapatkan akses ke service Google Sheets
	GetSheetsService(ctx context.Context) (*sheets.Service, error)

	// GetDriveService mendapatkan akses ke service Google Drive
	GetDriveService(ctx context.Context) (*drive.Service, error)

	// IsConfigured memeriksa apakah repository sudah dikonfigurasi
	IsConfigured() bool

	// GetCredentialsPath mendapatkan path ke file credentials
	GetCredentialsPath() string

	// SaveToken menyimpan token OAuth2
	SaveToken(token *oauth2.Token) error

	// GetAuthURL mendapatkan URL untuk otorisasi
	GetAuthURL(redirectURL string) (string, error)
}
