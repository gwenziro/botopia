package web

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gwenziro/botopia/internal/domain/finance"
	"github.com/gwenziro/botopia/internal/domain/service"
	"github.com/gwenziro/botopia/internal/usecase/command/list"
	"github.com/gwenziro/botopia/internal/usecase/stats"
)

// DashboardController controllers untuk dashboard web
type DashboardController struct {
	getStatsUseCase     *stats.GetStatsUseCase
	listCommandsUseCase *list.ListCommandsUseCase
	financeService      service.FinanceService
}

// NewDashboardController membuat instance controller baru
func NewDashboardController(
	statsUC *stats.GetStatsUseCase,
	cmdListUC *list.ListCommandsUseCase,
	financeService service.FinanceService,
) *DashboardController {
	return &DashboardController{
		getStatsUseCase:     statsUC,
		listCommandsUseCase: cmdListUC,
		financeService:      financeService,
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

	// Serialize commands untuk JavaScript
	commandsJSON, err := json.Marshal(commands)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to serialize commands: " + err.Error(),
		})
	}

	// Dapatkan transaksi terbaru (jika financeService tersedia)
	var recentTransactionsJSON string
	var weeklyStats map[string]interface{}
	if c.financeService != nil {
		// Buat context dengan timeout
		timeoutCtx, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
		defer cancel()

		// Ambil transaksi terbaru (5 transaksi)
		recentRecords, err := c.financeService.GetRecentRecords(timeoutCtx, 5)
		if err == nil && len(recentRecords) > 0 {
			// Konversi ke format yang sesuai untuk frontend
			transactions := make([]map[string]interface{}, 0, len(recentRecords))
			for _, record := range recentRecords {
				transactions = append(transactions, map[string]interface{}{
					"id":          record.UniqueCode,
					"date":        record.Date.Format("02 Jan 2006"),
					"description": record.Description,
					"amount":      record.Amount,
					"type":        string(record.Type),
					"category":    record.Category,
					"hasProof":    record.ProofURL != "" && record.ProofURL != "-",
				})
			}

			// Serialize transaksi terbaru
			if txJSON, err := json.Marshal(transactions); err == nil {
				recentTransactionsJSON = string(txJSON)
			}

			// Hitung statistik mingguan
			weeklyStats = calculateWeeklyStats(recentRecords)
		}
	}

	if recentTransactionsJSON == "" {
		recentTransactionsJSON = "[]" // Default empty array
	}

	// Render dashboard dengan data yang lengkap
	return ctx.Render("pages/dashboard", fiber.Map{
		"Title":               "Dashboard | Botopia",
		"Page":                "dashboard",
		"ConnectionState":     stats.ConnectionState,
		"IsConnected":         stats.IsConnected,
		"CommandCount":        stats.CommandCount,
		"MessageCount":        stats.MessageCount,
		"CommandsRun":         stats.CommandsRun,
		"CommandsJSON":        string(commandsJSON),
		"Phone":               stats.Phone,
		"Name":                stats.Name,
		"Uptime":              stats.Uptime,
		"RecentTransactions":  recentTransactionsJSON,
		"WeeklyStats":         weeklyStats,
		"SpreadsheetURL":      c.getSpreadsheetURL(),
		"HasFinanceService":   c.financeService != nil,
		"HasGoogleAPIService": c.checkGoogleAPIConfigured(),
		"SpreadsheetUrl":      c.getSpreadsheetURL(),
		"SpreadsheetId":       c.getSpreadsheetID(), // Add this
		"DriveFolderId":       c.getDriveFolderID(), // Already added
		"Commands":            commands,
		"ConnectivityUrl":     "/connectivity", // Tambahkan data untuk URL konektivitas
	}, "layouts/main")
}

// Fungsi helper untuk mendapatkan URL spreadsheet
func (c *DashboardController) getSpreadsheetURL() string {
	if c.financeService != nil {
		return c.financeService.GetSpreadsheetURL()
	}
	return "#"
}

// Fungsi helper untuk memeriksa apakah Google API dikonfigurasi
func (c *DashboardController) checkGoogleAPIConfigured() bool {
	return c.financeService != nil && c.financeService.GetSpreadsheetURL() != "#"
}

// calculateWeeklyStats menghitung statistik keuangan mingguan
func calculateWeeklyStats(records []*finance.FinanceRecord) map[string]interface{} {
	// Inisialisasi statistik
	totalIncome := 0.0
	totalExpense := 0.0
	categories := make(map[string]float64)

	// Hitung hanya untuk data dari 7 hari terakhir
	oneWeekAgo := time.Now().AddDate(0, 0, -7)

	for _, record := range records {
		// Lewati jika record lebih dari 7 hari yang lalu
		if record.Date.Before(oneWeekAgo) {
			continue
		}

		if record.Type == finance.TypeIncome {
			totalIncome += record.Amount
		} else {
			totalExpense += record.Amount
			categories[record.Category] += record.Amount
		}
	}

	// Temukan kategori pengeluaran terbesar
	largestCategory := ""
	largestAmount := 0.0
	for cat, amount := range categories {
		if amount > largestAmount {
			largestAmount = amount
			largestCategory = cat
		}
	}

	return map[string]interface{}{
		"totalIncome":     totalIncome,
		"totalExpense":    totalExpense,
		"balance":         totalIncome - totalExpense,
		"largestCategory": largestCategory,
		"largestAmount":   largestAmount,
		"categoryData":    categories,
	}
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

// HandleGetRecentTransactions mengambil transaksi terbaru
func (c *DashboardController) HandleGetRecentTransactions(ctx *fiber.Ctx) error {
	if c.financeService == nil {
		return ctx.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Finance service not available",
		})
	}

	// Buat context dengan timeout
	timeoutCtx, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	// Parameter limit dari query string
	limit := 5 // default
	if ctx.QueryInt("limit", 0) > 0 {
		limit = ctx.QueryInt("limit", 5)
	}

	// Ambil transaksi terbaru
	recentRecords, err := c.financeService.GetRecentRecords(timeoutCtx, limit)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get recent transactions: " + err.Error(),
		})
	}

	// Konversi ke format yang sesuai untuk frontend
	transactions := make([]map[string]interface{}, 0, len(recentRecords))
	for _, record := range recentRecords {
		transactions = append(transactions, map[string]interface{}{
			"id":          record.UniqueCode,
			"date":        record.Date.Format("02 Jan 2006"),
			"description": record.Description,
			"amount":      record.Amount,
			"type":        string(record.Type),
			"category":    record.Category,
			"hasProof":    record.ProofURL != "" && record.ProofURL != "-",
		})
	}

	return ctx.JSON(fiber.Map{
		"transactions": transactions,
	})
}

// Fungsi helper untuk mendapatkan Drive Folder ID
func (c *DashboardController) getDriveFolderID() string {
	if c.financeService == nil {
		return ""
	}

	// Implementasi sebenarnya akan bergantung pada struktur aplikasi
	// Ini adalah contoh sederhana
	if service, ok := c.financeService.(interface{ GetDriveFolderID() string }); ok {
		return service.GetDriveFolderID()
	}

	// Fallback jika interface tidak tersedia
	return ""
}

// getSpreadsheetID mendapatkan ID spreadsheet
func (c *DashboardController) getSpreadsheetID() string {
	if c.financeService == nil {
		return ""
	}

	return c.financeService.GetSpreadsheetID()
}
