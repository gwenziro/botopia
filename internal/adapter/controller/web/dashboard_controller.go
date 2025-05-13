package web

import (
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

	// Render dashboard
	return ctx.Render("pages/dashboard", fiber.Map{
		"Title":           "Dashboard | Botopia",
		"Page":            "dashboard",
		"ConnectionState": stats.ConnectionState,
		"CommandCount":    stats.CommandCount,
		"MessageCount":    stats.MessageCount,
		"CommandsRun":     stats.CommandsRun,
		"Commands":        commands,
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
