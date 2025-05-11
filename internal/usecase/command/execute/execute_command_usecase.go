package execute

import (
	"context"
	"errors"

	"github.com/gwenziro/botopia/internal/domain/message"
	"github.com/gwenziro/botopia/internal/domain/repository"
)

// ExecuteCommandUseCase menangani eksekusi command dari pesan
type ExecuteCommandUseCase struct {
	commandRepo   repository.CommandRepository
	statsRepo     repository.StatsRepository
	connRepo      repository.ConnectionRepository
	commandPrefix string
}

// NewExecuteCommandUseCase membuat instance usecase baru
func NewExecuteCommandUseCase(
	cmdRepo repository.CommandRepository,
	statsRepo repository.StatsRepository,
	connRepo repository.ConnectionRepository,
	prefix string,
) *ExecuteCommandUseCase {
	return &ExecuteCommandUseCase{
		commandRepo:   cmdRepo,
		statsRepo:     statsRepo,
		connRepo:      connRepo,
		commandPrefix: prefix,
	}
}

// ErrCommandNotFound adalah error ketika command tidak ditemukan
var ErrCommandNotFound = errors.New("command tidak ditemukan")

// Execute menjalankan command dari pesan
func (uc *ExecuteCommandUseCase) Execute(ctx context.Context, msg *message.Message) (string, error) {
	// Validasi input
	if msg == nil {
		return "", errors.New("message cannot be nil")
	}

	// Ekstrak command dari pesan
	cmdName, args, isCommand := msg.ExtractCommand(uc.commandPrefix)
	if !isCommand {
		return "", nil // Bukan command, abaikan
	}

	// Cari command berdasarkan nama
	cmd, found := uc.commandRepo.FindByName(cmdName)
	if !found {
		return "", ErrCommandNotFound
	}

	// Eksekusi command
	response, err := cmd.Execute(args, msg)
	if err != nil {
		return "", err
	}

	// Update statistik
	uc.statsRepo.IncrementCommandCount()

	return response, nil
}
