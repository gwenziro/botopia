package connection

import (
	"context"
	"errors"
	"sync"

	"github.com/gwenziro/botopia/internal/domain/repository"
	"github.com/gwenziro/botopia/internal/domain/user"
	"github.com/gwenziro/botopia/internal/usecase/dto"
)

// ConnectWhatsAppUseCase mengatur koneksi WhatsApp
type ConnectWhatsAppUseCase struct {
	connectionRepo repository.ConnectionRepository
}

// NewConnectWhatsAppUseCase membuat use case baru
func NewConnectWhatsAppUseCase(repo repository.ConnectionRepository) *ConnectWhatsAppUseCase {
	return &ConnectWhatsAppUseCase{
		connectionRepo: repo,
	}
}

// Tambahkan variabel state dan mutex
var (
	currentQRCode  string
	isConnecting   bool
	connectionLock sync.Mutex
)

// IsConnecting memeriksa apakah sedang dalam proses connecting
func (uc *ConnectWhatsAppUseCase) IsConnecting() bool {
	connectionLock.Lock()
	defer connectionLock.Unlock()
	return isConnecting
}

// Execute menjalankan koneksi ke WhatsApp
// Mengembalikan status koneksi dan error jika ada
func (uc *ConnectWhatsAppUseCase) Execute(ctx context.Context) (*dto.ConnectionStatusDTO, error) {
	if uc.connectionRepo.IsConnected() {
		return &dto.ConnectionStatusDTO{
			IsConnected: true,
			Message:     "WhatsApp already connected",
		}, nil
	}

	// Set state connecting
	connectionLock.Lock()
	isConnecting = true
	connectionLock.Unlock()

	// Reset state when done
	defer func() {
		connectionLock.Lock()
		isConnecting = false
		connectionLock.Unlock()
	}()

	err := uc.connectionRepo.Connect(ctx)
	if err != nil {
		return &dto.ConnectionStatusDTO{
			IsConnected: false,
			Message:     "Failed to connect: " + err.Error(),
			Error:       err,
		}, err
	}

	// Dapatkan informasi user
	user, _ := uc.connectionRepo.GetCurrentUser()

	// Buat response
	response := &dto.ConnectionStatusDTO{
		IsConnected: uc.connectionRepo.IsConnected(),
		Message:     "Successfully connected to WhatsApp",
	}

	if user != nil {
		response.Phone = user.Phone
	}

	return response, nil
}

// Disconnect memutuskan koneksi WhatsApp
func (uc *ConnectWhatsAppUseCase) Disconnect() error {
	if !uc.connectionRepo.IsConnected() {
		return errors.New("not connected")
	}

	return uc.connectionRepo.Disconnect()
}

// GetQRChannel mendapatkan channel QR code
func (uc *ConnectWhatsAppUseCase) GetQRChannel() <-chan string {
	// Wrap channel asli dengan channel buffered untuk menghindari blocking
	originalChan := uc.connectionRepo.GetQRChannel()
	bufferedChan := make(chan string, 1)

	// Goroutine untuk memproses QR dari channel asli
	go func() {
		for qr := range originalChan {
			// Simpan QR code terbaru
			connectionLock.Lock()
			currentQRCode = qr
			connectionLock.Unlock()

			// Forward ke channel buffer, non-blocking
			select {
			case bufferedChan <- qr:
				// QR berhasil dikirim ke channel
			default:
				// Channel penuh, abaikan (non-blocking)
			}
		}
	}()

	return bufferedChan
}

// GetCurrentQR returns the latest QR code (if available)
func (uc *ConnectWhatsAppUseCase) GetCurrentQR() string {
	connectionLock.Lock()
	defer connectionLock.Unlock()
	return currentQRCode
}

// IsConnected memeriksa status koneksi
func (uc *ConnectWhatsAppUseCase) IsConnected() bool {
	return uc.connectionRepo.IsConnected()
}

// GetCurrentUser mendapatkan informasi user yang terhubung
func (uc *ConnectWhatsAppUseCase) GetCurrentUser() (*user.User, error) {
	return uc.connectionRepo.GetCurrentUser()
}
