package stats

import (
	"context"

	"github.com/gwenziro/botopia/internal/domain/repository"
	"github.com/gwenziro/botopia/internal/usecase/dto"
)

// GetStatsUseCase use case untuk mendapatkan statistik bot
type GetStatsUseCase struct {
	statsRepo      repository.StatsRepository
	connectionRepo repository.ConnectionRepository
	commandRepo    repository.CommandRepository
}

// NewGetStatsUseCase membuat instance use case baru
func NewGetStatsUseCase(
	statsRepo repository.StatsRepository,
	connectionRepo repository.ConnectionRepository,
	commandRepo repository.CommandRepository,
) *GetStatsUseCase {
	return &GetStatsUseCase{
		statsRepo:      statsRepo,
		connectionRepo: connectionRepo,
		commandRepo:    commandRepo,
	}
}

// Execute menjalankan pengambilan statistik
func (uc *GetStatsUseCase) Execute(ctx context.Context) (*dto.StatsDTO, error) {
	// Dapatkan statistik dari repository
	stats, err := uc.statsRepo.GetStats()
	if err != nil {
		return nil, err
	}

	// Dapatkan status koneksi
	isConnected := uc.connectionRepo.IsConnected()

	// Dapatkan info user jika terhubung
	var phone string
	if isConnected {
		user, _ := uc.connectionRepo.GetCurrentUser()
		if user != nil {
			phone = user.Phone
		}
	}

	// Dapatkan jumlah command tersedia
	commands := uc.commandRepo.GetAll()
	commandCount := len(commands)

	// Buat DTO
	return &dto.StatsDTO{
		ConnectionState: stats.ConnectionState,
		IsConnected:     isConnected,
		MessageCount:    stats.MessageCount,
		CommandsRun:     stats.CommandsRun,
		Uptime:          stats.Uptime,
		CommandCount:    commandCount,
		Phone:           phone,
		Name:            "WhatsApp User", // Nilai default
	}, nil
}
