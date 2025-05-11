package list

import (
	"context"

	"github.com/gwenziro/botopia/internal/domain/repository"
	"github.com/gwenziro/botopia/internal/usecase/dto"
)

// ListCommandsUseCase use case untuk mendaftar command
type ListCommandsUseCase struct {
	commandRepo repository.CommandRepository
}

// NewListCommandsUseCase membuat instance use case baru
func NewListCommandsUseCase(commandRepo repository.CommandRepository) *ListCommandsUseCase {
	return &ListCommandsUseCase{
		commandRepo: commandRepo,
	}
}

// Execute menjalankan pendaftaran command
func (uc *ListCommandsUseCase) Execute(ctx context.Context) (map[string]dto.CommandDTO, error) {
	// Ambil semua command
	commands := uc.commandRepo.GetAll()

	// Konversi ke DTO
	result := make(map[string]dto.CommandDTO, len(commands))
	for name, cmd := range commands {
		result[name] = dto.CommandDTO{
			Name:        cmd.GetName(),
			Description: cmd.GetDescription(),
		}
	}

	return result, nil
}
