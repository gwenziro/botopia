package web

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
	"github.com/gwenziro/botopia/internal/usecase/command/list"
)

// CommandsController menangani request terkait daftar commands
type CommandsController struct {
	listCommandsUseCase *list.ListCommandsUseCase
	log                 *logger.Logger
}

// NewCommandsController membuat instance controller baru
func NewCommandsController(listCommandsUseCase *list.ListCommandsUseCase) *CommandsController {
	return &CommandsController{
		listCommandsUseCase: listCommandsUseCase,
		log:                 logger.New("CommandsController", logger.INFO, true),
	}
}

// HandleCommandsPage menangani halaman commands
func (c *CommandsController) HandleCommandsPage(ctx *fiber.Ctx) error {
	// Buat context dengan timeout
	timeoutCtx, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	// Dapatkan daftar commands
	commands, err := c.listCommandsUseCase.Execute(timeoutCtx)
	if err != nil {
		c.log.Error("Failed to get commands: %v", err)
		return ctx.Render("pages/commands", fiber.Map{
			"Title":        "Commands | Botopia",
			"Page":         "commands",
			"CommandsJSON": "{}",
		}, "layouts/main")
	}

	// Serialisasi commands ke JSON untuk digunakan di frontend
	commandsJSON, err := json.Marshal(commands)
	if err != nil {
		c.log.Error("Failed to serialize commands: %v", err)
		return ctx.Render("pages/commands", fiber.Map{
			"Title":        "Commands | Botopia",
			"Page":         "commands",
			"CommandsJSON": "{}",
		}, "layouts/main")
	}

	// Render halaman
	return ctx.Render("pages/commands", fiber.Map{
		"Title":        "Commands | Botopia",
		"Page":         "commands",
		"CommandsJSON": string(commandsJSON),
	}, "layouts/main")
}
