package command

import (
	"github.com/gwenziro/botopia/internal/domain/command/finance"
	"github.com/gwenziro/botopia/internal/domain/command/help"
	"github.com/gwenziro/botopia/internal/domain/command/ping"
	"github.com/gwenziro/botopia/internal/domain/repository"
	"github.com/gwenziro/botopia/internal/domain/service"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
)

// CommandInitializer mendaftarkan command-command default
type CommandInitializer struct {
	cmdRepo        repository.CommandRepository
	financeService service.FinanceService
	log            *logger.Logger
}

// NewCommandInitializer membuat instance initializer baru
func NewCommandInitializer(
	cmdRepo repository.CommandRepository,
	financeService service.FinanceService,
) *CommandInitializer {
	return &CommandInitializer{
		cmdRepo:        cmdRepo,
		financeService: financeService,
		log:            logger.New("CommandInitializer", logger.INFO, true),
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
	helpCmd := help.NewCommand(c.cmdRepo) // Gunakan help bukan guide
	c.cmdRepo.Register(helpCmd)
	c.log.Info("Command '%s' terdaftar", helpCmd.GetName())

	// 3. Finance commands
	if c.financeService != nil {
		// Pengeluaran command
		expenseCmd := finance.NewAddExpenseCommand(c.financeService)
		c.cmdRepo.Register(expenseCmd)
		c.log.Info("Command '%s' terdaftar", expenseCmd.GetName())

		// Pemasukan command
		incomeCmd := finance.NewAddIncomeCommand(c.financeService)
		c.cmdRepo.Register(incomeCmd)
		c.log.Info("Command '%s' terdaftar", incomeCmd.GetName())

		// Upload bukti transaksi command
		uploadCmd := finance.NewUploadProofCommand(c.financeService)
		c.cmdRepo.Register(uploadCmd)
		c.log.Info("Command '%s' terdaftar", uploadCmd.GetName())
	} else {
		c.log.Warn("Finance service tidak tersedia, command finance tidak akan didaftarkan")
	}

	// Tambahkan command default lainnya di sini
}

// GetCommandCount mendapatkan jumlah command yang terdaftar
func (c *CommandInitializer) GetCommandCount() int {
	return len(c.cmdRepo.GetAll())
}
