package help

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gwenziro/botopia/internal/domain/command/common"
	"github.com/gwenziro/botopia/internal/domain/message"
	"github.com/gwenziro/botopia/internal/domain/repository"
)

// Command implementasi command help
type Command struct {
	common.BaseCommand
	cmdRepo repository.CommandRepository
}

// NewCommand membuat instance help command baru
func NewCommand(cmdRepo repository.CommandRepository) *Command {
	cmd := &Command{
		cmdRepo: cmdRepo,
	}
	cmd.Name = "panduan"
	cmd.Description = "Menampilkan daftar command yang tersedia"
	cmd.Category = "Help"
	cmd.Usage = "!panduan [nama_command]"
	return cmd
}

// Execute menjalankan command panduan
func (c *Command) Execute(args []string, msg *message.Message) (string, error) {
	if len(args) > 0 {
		// Tampilkan bantuan untuk command tertentu
		cmdName := args[0]
		cmd, found := c.cmdRepo.FindByName(cmdName)
		if !found {
			return fmt.Sprintf("Command '%s' tidak ditemukan", cmdName), nil
		}

		// Tampilkan informasi lengkap tentang command
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("*%s*: %s\n\n", cmd.GetName(), cmd.GetDescription()))

		// Tampilkan kategori jika tersedia
		category := cmd.GetCategory()
		if category != "" {
			sb.WriteString(fmt.Sprintf("Kategori: %s\n", category))
		}

		// Tampilkan penggunaan yang benar
		usage := cmd.GetUsage()
		if usage != "" {
			sb.WriteString(fmt.Sprintf("Penggunaan: %s\n", usage))
		} else {
			sb.WriteString(fmt.Sprintf("Penggunaan: !%s\n", cmd.GetName()))
		}

		return sb.String(), nil
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

	sb.WriteString("\nGunakan !panduan <command> untuk informasi lebih lanjut.")
	return sb.String(), nil
}
