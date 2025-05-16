package web

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gwenziro/botopia/internal/domain/service"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
)

// DataMasterController adalah controller untuk pengelolaan data master
type DataMasterController struct {
	financeService service.FinanceService
	log            *logger.Logger
}

// NewDataMasterController membuat instance controller baru
func NewDataMasterController(financeService service.FinanceService) *DataMasterController {
	return &DataMasterController{
		financeService: financeService,
		log:            logger.New("DataMasterController", logger.INFO, true),
	}
}

// HandleDataMasterPage menangani halaman data master
func (c *DataMasterController) HandleDataMasterPage(ctx *fiber.Ctx) error {
	// Buat context dengan timeout
	timeoutCtx, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	// Dapatkan konfigurasi untuk menampilkan data master
	config, err := c.financeService.GetConfiguration(timeoutCtx)
	if err != nil {
		c.log.Error("Gagal memuat konfigurasi: %v", err)
		// Tetap render halaman, tetapi dengan data kosong
		return ctx.Render("pages/data_master", fiber.Map{
			"Title":      "Data Master | Botopia",
			"Page":       "data-master",
			"ActiveTab":  ctx.Query("tab", "expense-categories"),
			"Error":      "Gagal memuat data: " + err.Error(),
			"ConfigJSON": "{}",
		}, "layouts/main")
	}

	// Pastikan slice tidak nil untuk menghindari masalah di frontend
	if config.ExpenseCategories == nil {
		config.ExpenseCategories = []string{}
	}
	if config.IncomeCategories == nil {
		config.IncomeCategories = []string{}
	}
	if config.PaymentMethods == nil {
		config.PaymentMethods = []string{}
	}
	if config.StorageMedias == nil {
		config.StorageMedias = []string{}
	}

	// Buat data untuk ditampilkan di halaman
	data := fiber.Map{
		"Title":             "Data Master | Botopia",
		"Page":              "data-master",
		"ExpenseCategories": config.ExpenseCategories,
		"IncomeCategories":  config.IncomeCategories,
		"PaymentMethods":    config.PaymentMethods,
		"StorageMedias":     config.StorageMedias,
		"ActiveTab":         ctx.Query("tab", "expense-categories"),
	}

	// Serialize konfigurasi untuk JavaScript
	configData := map[string]interface{}{
		"expenseCategories": config.ExpenseCategories, // lowercase keys
		"incomeCategories":  config.IncomeCategories,
		"paymentMethods":    config.PaymentMethods,
		"storageMedias":     config.StorageMedias,
	}

	configJSON, err := json.Marshal(configData)
	if err != nil {
		c.log.Error("Gagal menyerialisasi konfigurasi: %v", err)
		data["ConfigJSON"] = "{}"
	} else {
		c.log.Info("JSON data serialized successfully, length: %d bytes", len(configJSON))
		data["ConfigJSON"] = string(configJSON)
	}

	// Render halaman data master
	return ctx.Render("pages/data_master", data, "layouts/main")
}

// HandleGetMasterData menangani API untuk mendapatkan data master
func (c *DataMasterController) HandleGetMasterData(ctx *fiber.Ctx) error {
	// Buat context dengan timeout
	timeoutCtx, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	// Dapatkan konfigurasi
	config, err := c.financeService.GetConfiguration(timeoutCtx)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal memuat konfigurasi: " + err.Error(),
		})
	}

	// Kembalikan data yang diminta berdasarkan parameter
	dataType := ctx.Query("type", "all")

	switch dataType {
	case "expense-categories":
		return ctx.JSON(fiber.Map{"data": config.ExpenseCategories})
	case "income-categories":
		return ctx.JSON(fiber.Map{"data": config.IncomeCategories})
	case "payment-methods":
		return ctx.JSON(fiber.Map{"data": config.PaymentMethods})
	case "storage-medias":
		return ctx.JSON(fiber.Map{"data": config.StorageMedias})
	default:
		return ctx.JSON(fiber.Map{
			"expenseCategories": config.ExpenseCategories,
			"incomeCategories":  config.IncomeCategories,
			"paymentMethods":    config.PaymentMethods,
			"storageMedias":     config.StorageMedias,
		})
	}
}
