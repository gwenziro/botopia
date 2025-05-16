package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gwenziro/botopia/internal/usecase/command/list"
	"github.com/gwenziro/botopia/internal/usecase/stats"
)

// DashboardController controllers untuk dashboard web
type DashboardController struct {
	getStatsUseCase     *stats.GetStatsUseCase
	listCommandsUseCase *list.ListCommandsUseCase
}

// NewDashboardController membuat instance controller baru
func NewDashboardController(
	statsUC *stats.GetStatsUseCase,
	cmdListUC *list.ListCommandsUseCase,
) *DashboardController {
	return &DashboardController{
		getStatsUseCase:     statsUC,
		listCommandsUseCase: cmdListUC,
	}
}

// HandleIndex menangani halaman index
func (c *DashboardController) HandleIndex(ctx *fiber.Ctx) error {
	return ctx.Render("pages/index", fiber.Map{
		"Title": "Botopia - WhatsApp Bot untuk Keuangan",
		"Page":  "home",
	}, "layouts/home_layout")
}

// HandleDashboard menangani halaman dashboard
func (c *DashboardController) HandleDashboard(ctx *fiber.Ctx) error {
	// Dapatkan statistik
	stats, err := c.getStatsUseCase.Execute(ctx.Context())
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get stats: " + err.Error(),
		})
	}

	// Dapatkan command
	commands, err := c.listCommandsUseCase.Execute(ctx.Context())
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list commands: " + err.Error(),
		})
	}

	// Debugging: cetak commands untuk melihat strukturnya
	for name, cmd := range commands {
		fmt.Printf("Command: %s, Description: %s\n", name, cmd.Description)
	}

	// Serialize commands dengan format yang pasti cocok dengan yang diperlukan oleh Alpine.js
	commandMap := make(map[string]map[string]interface{})
	for name, cmd := range commands {
		commandMap[name] = map[string]interface{}{
			"description": cmd.Description,
			"name":        name,
		}
	}

	commandsJSON, err := json.Marshal(commandMap)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to serialize commands: " + err.Error(),
		})
	}

	// Render dashboard dengan data yang lengkap
	return ctx.Render("pages/dashboard", fiber.Map{
		"Title":           "Dashboard | Botopia",
		"Page":            "dashboard",
		"ConnectionState": stats.ConnectionState,
		"IsConnected":     stats.IsConnected,
		"CommandCount":    stats.CommandCount,
		"MessageCount":    stats.MessageCount,
		"CommandsRun":     stats.CommandsRun,
		"CommandsJSON":    string(commandsJSON),
		"Phone":           stats.Phone,
		"Name":            stats.Name,
	}, "layouts/main")
}

// HandleGetStats menangani API stats
func (c *DashboardController) HandleGetStats(ctx *fiber.Ctx) error {
	// Dapatkan statistik
	stats, err := c.getStatsUseCase.Execute(ctx.Context())
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get stats: " + err.Error(),
		})
	}

	// Return JSON
	return ctx.JSON(stats)
}

// HandleGetCommands menangani API untuk mendapatkan daftar command
func (c *DashboardController) HandleGetCommands(ctx *fiber.Ctx) error {
	// Dapatkan command
	commands, err := c.listCommandsUseCase.Execute(ctx.Context())
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list commands: " + err.Error(),
		})
	}

	// Return JSON
	return ctx.JSON(commands)
}
