package command

import (
	"github.com/gwenziro/botopia/internal/domain/command/help"
	"github.com/gwenziro/botopia/internal/domain/command/ping"
	"github.com/gwenziro/botopia/internal/domain/repository"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
)

// CommandInitializer mendaftarkan command-command default
type CommandInitializer struct {
	cmdRepo repository.CommandRepository
	log     *logger.Logger
}

// NewCommandInitializer membuat instance initializer baru
func NewCommandInitializer(cmdRepo repository.CommandRepository) *CommandInitializer {
	return &CommandInitializer{
		cmdRepo: cmdRepo,
		log:     logger.New("CommandInitializer", logger.INFO, true),
	}
}

// RegisterDefaultCommands mendaftarkan semua command default
func (c *CommandInitializer) RegisterDefaultCommands() {
	c.log.Info("Mendaftarkan command default")

	// 1. Ping command
	pingCmd := ping.NewCommand()
	c.cmdRepo.Register(pingCmd)
	c.log.Info("Command '%s' terdaftar", pingCmd.GetName())

	// 2. Help command
	helpCmd := help.NewCommand(c.cmdRepo)
	c.cmdRepo.Register(helpCmd)
	c.log.Info("Command '%s' terdaftar", helpCmd.GetName())

	// Tambahkan command default lainnya di sini
}

// GetCommandCount mendapatkan jumlah command yang terdaftar
func (c *CommandInitializer) GetCommandCount() int {
	return len(c.cmdRepo.GetAll())
}
