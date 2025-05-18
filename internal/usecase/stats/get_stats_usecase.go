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
	name := "WhatsApp User" // Nilai default
	var deviceDetailsDTO *dto.DeviceDetailsDTO

	if isConnected {
		user, err := uc.connectionRepo.GetCurrentUser()
		if err == nil && user != nil {
			phone = user.Phone

			// Prioritaskan PushName, lalu Name
			if user.PushName != "" {
				name = user.PushName
			} else if user.Name != "" {
				name = user.Name
			}

			// Tambahkan detail perangkat jika tersedia
			if user.DeviceDetails != nil {
				deviceDetailsDTO = &dto.DeviceDetailsDTO{
					Platform:    user.DeviceDetails.Platform,
					DeviceModel: user.DeviceDetails.DeviceModel,
					OSVersion:   user.DeviceDetails.OSVersion,
					ClientType:  user.DeviceDetails.ClientType,
					IPAddress:   user.DeviceDetails.IPAddress,
					DeviceID:    user.DeviceDetails.DeviceID,
				}
			}
		}
	}

	// Dapatkan jumlah command tersedia
	commands := uc.commandRepo.GetAll()
	commandCount := len(commands)

	// Buat DTO
	statsDTO := &dto.StatsDTO{
		ConnectionState: stats.ConnectionState,
		IsConnected:     isConnected,
		MessageCount:    stats.MessageCount,
		CommandsRun:     stats.CommandsRun,
		Uptime:          stats.Uptime,
		SystemUptime:    stats.SystemUptime,
		CommandCount:    commandCount,
		Phone:           phone,
		Name:            name,
		DeviceDetails:   deviceDetailsDTO,
	}

	return statsDTO, nil
}
