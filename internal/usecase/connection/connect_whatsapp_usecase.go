package connection

import (
	"context"
	"errors"

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

// Execute menjalankan koneksi ke WhatsApp
// Mengembalikan status koneksi dan error jika ada
func (uc *ConnectWhatsAppUseCase) Execute(ctx context.Context) (*dto.ConnectionStatusDTO, error) {
	if uc.connectionRepo.IsConnected() {
		return &dto.ConnectionStatusDTO{
			IsConnected: true,
			Message:     "WhatsApp already connected",
		}, nil
	}

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
	return uc.connectionRepo.GetQRChannel()
}

// IsConnected memeriksa status koneksi
func (uc *ConnectWhatsAppUseCase) IsConnected() bool {
	return uc.connectionRepo.IsConnected()
}

// GetCurrentUser mendapatkan informasi user yang terhubung
func (uc *ConnectWhatsAppUseCase) GetCurrentUser() (*user.User, error) {
	return uc.connectionRepo.GetCurrentUser()
}
