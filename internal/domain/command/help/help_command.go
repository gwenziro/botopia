package help

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gwenziro/botopia/internal/domain/message"
	"github.com/gwenziro/botopia/internal/domain/repository"
)

// Command implementasi command help
type Command struct {
	cmdRepo repository.CommandRepository
}

// NewCommand membuat instance help command baru
func NewCommand(cmdRepo repository.CommandRepository) *Command {
	return &Command{
		cmdRepo: cmdRepo,
	}
}

// GetName mengembalikan nama command
func (c *Command) GetName() string {
	return "help"
}

// GetDescription mengembalikan deskripsi command
func (c *Command) GetDescription() string {
	return "Menampilkan daftar command yang tersedia"
}

// Execute menjalankan command help
func (c *Command) Execute(args []string, msg *message.Message) (string, error) {
	if len(args) > 0 {
		// Tampilkan bantuan untuk command tertentu
		cmdName := args[0]
		cmd, found := c.cmdRepo.FindByName(cmdName)
		if !found {
			return fmt.Sprintf("Command '%s' tidak ditemukan", cmdName), nil
		}

		return fmt.Sprintf("*%s*: %s", cmd.GetName(), cmd.GetDescription()), nil
	}

	// Tampilkan semua command yang tersedia
	commands := c.cmdRepo.GetAll()

	// Ekstraksi & urutkan nama command
	var commandNames []string
	for name := range commands {
		commandNames = append(commandNames, name)
	}
	sort.Strings(commandNames)

	var sb strings.Builder
	sb.WriteString("*Command yang tersedia:*\n\n")

	for _, name := range commandNames {
		cmd := commands[name]
		sb.WriteString(fmt.Sprintf("!%s - %s\n", name, cmd.GetDescription()))
	}

	sb.WriteString("\nGunakan !help <command> untuk informasi lebih lanjut.")
	return sb.String(), nil
}
