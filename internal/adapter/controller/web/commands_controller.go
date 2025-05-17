package web

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gwenziro/botopia/internal/domain/command"
	"github.com/gwenziro/botopia/internal/domain/repository"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
)

// CommandsController adalah controller untuk halaman commands
type CommandsController struct {
	commandRepo repository.CommandRepository
	log         *logger.Logger
}

// NewCommandsController membuat instance commands controller baru
func NewCommandsController(commandRepo repository.CommandRepository) *CommandsController {
	return &CommandsController{
		commandRepo: commandRepo,
		log:         logger.New("CommandsController", logger.INFO, true),
	}
}

// HandleCommandsPage menangani halaman daftar commands
func (c *CommandsController) HandleCommandsPage(ctx *fiber.Ctx) error {
	// Buat context dengan timeout
	_, cancel := context.WithTimeout(ctx.Context(), 5*time.Second)
	defer cancel()

	// Ambil semua command yang tersedia
	commands := c.commandRepo.GetAll()

	// Ekstraksi nama-nama command dan urutkan
	var commandNames []string
	for name := range commands {
		commandNames = append(commandNames, name)
	}
	sort.Strings(commandNames)

	// Ekstrak kategori unik (jika tersedia)
	categories := make(map[string]bool)
	for _, cmd := range commands {
		// Jika command memiliki method GetCategory, gunakan
		if categorizer, ok := cmd.(interface{ GetCategory() string }); ok {
			category := categorizer.GetCategory()
			if category != "" {
				categories[category] = true
			}
		}
	}

	// Konversi map kategori ke slice
	var uniqueCategories []string
	for category := range categories {
		uniqueCategories = append(uniqueCategories, category)
	}
	sort.Strings(uniqueCategories)

	// Debugging: log jumlah command
	c.log.Info("Total commands: %d, categories: %d", len(commands), len(uniqueCategories))

	// Buat data untuk template
	data := fiber.Map{
		"Title":          "Commands | Botopia",
		"Page":           "commands",
		"Commands":       commands,
		"Categories":     uniqueCategories,
		"Total":          len(commands),
		"SearchQuery":    ctx.Query("q", ""),
		"CategoryFilter": ctx.Query("category", ""),
	}

	// Serialisasi commands untuk JavaScript - PERBAIKAN: Langsung pake hasil tanpa string tambahan
	commandsJSON, err := json.Marshal(c.getCommandsForAPI(commands))
	if err != nil {
		c.log.Error("Failed to serialize commands: %v", err)
		data["CommandsJSON"] = "{}"
	} else {
		c.log.Info("Commands serialized successfully, JSON length: %d bytes", len(commandsJSON))
		// Langsung gunakan hasil serialisasi tanpa serialisasi tambahan
		data["CommandsJSON"] = string(commandsJSON)
	}

	return ctx.Render("pages/commands", data, "layouts/main")
}

// HandleGetCommands menangani API untuk mendapatkan daftar commands
func (c *CommandsController) HandleGetCommands(ctx *fiber.Ctx) error {
	// Buat context dengan timeout
	_, cancel := context.WithTimeout(ctx.Context(), 5*time.Second)
	defer cancel()

	// Ambil semua command
	commands := c.commandRepo.GetAll()

	// Format data untuk API
	result := c.getCommandsForAPI(commands)

	return ctx.Status(http.StatusOK).JSON(result)
}

// getCommandsForAPI mengonversi command ke format yang sesuai untuk API
func (c *CommandsController) getCommandsForAPI(cmds map[string]command.Command) map[string]map[string]interface{} {
	result := make(map[string]map[string]interface{})

	for name, cmd := range cmds {
		cmdData := map[string]interface{}{
			"name":        name,
			"description": cmd.GetDescription(),
		}

		// Tambahkan kategori jika tersedia
		category := "Umum" // Default category
		if getter, ok := cmd.(interface{ GetCategory() string }); ok {
			if cat := getter.GetCategory(); cat != "" {
				category = cat
			}
		}
		cmdData["category"] = category

		// Tambahkan usage jika tersedia
		usage := "!" + name // Default usage
		if getter, ok := cmd.(interface{ GetUsage() string }); ok {
			if u := getter.GetUsage(); u != "" {
				usage = u
			}
		}
		cmdData["usage"] = usage

		// Log setiap command untuk debug
		c.log.Debug("Command: %s, desc: %s, category: %s, usage: %s",
			name, cmdData["description"], cmdData["category"], cmdData["usage"])

		result[name] = cmdData
	}

	c.log.Info("Prepared %d commands for API response", len(result))
	return result
}
