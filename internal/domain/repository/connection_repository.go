package repository

import (
	"context"

	"github.com/gwenziro/botopia/internal/domain/message"
	"github.com/gwenziro/botopia/internal/domain/user"
)

// ConnectionRepository mendefinisikan kontrak untuk repository koneksi WhatsApp
type ConnectionRepository interface {
	// Connect menghubungkan ke WhatsApp
	Connect(ctx context.Context) error

	// Disconnect memutuskan koneksi WhatsApp
	Disconnect() error

	// IsConnected memeriksa status koneksi
	IsConnected() bool

	// GetCurrentUser mendapatkan informasi user yang terhubung
	GetCurrentUser() (*user.User, error)

	// SendMessage mengirim pesan
	SendMessage(ctx context.Context, chatID string, text string) error

	// RegisterMessageHandler mendaftarkan handler untuk pesan masuk
	RegisterMessageHandler(handler func(*message.Message))

	// GetQRChannel mendapatkan channel untuk QR code
	GetQRChannel() <-chan string

	// DownloadMedia mengunduh media dari pesan
	DownloadMedia(ctx context.Context, msg *message.Message) (string, error)
}
