package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gwenziro/botopia/internal/infrastructure/config"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
)

// ConfigController adalah controller untuk konfigurasi sistem
type ConfigController struct {
	config *config.Config
	log    *logger.Logger
}

// NewConfigController membuat instance controller konfigurasi baru
func NewConfigController(config *config.Config) *ConfigController {
	return &ConfigController{
		config: config,
		log:    logger.New("ConfigController", logger.INFO, true),
	}
}

// HandleConfigPage menangani halaman konfigurasi
func (c *ConfigController) HandleConfigPage(ctx *fiber.Ctx) error {
	return ctx.Render("pages/configuration", fiber.Map{
		"Title": "Konfigurasi | Botopia",
		"Page":  "konfigurasi",
	}, "layouts/main")
}

// HandleGetConfig menangani API untuk mendapatkan konfigurasi
func (c *ConfigController) HandleGetConfig(ctx *fiber.Ctx) error {
	// Kembalikan konfigurasi yang aman untuk frontend
	return ctx.JSON(fiber.Map{
		"commandPrefix":   c.config.CommandPrefix,
		"logLevel":        c.config.LogLevel,
		"webPort":         c.config.WebPort,
		"webHost":         c.config.WebHost,
		"spreadsheetId":   c.config.GoogleSheets.SpreadsheetID,
		"driveFolderId":   c.config.GoogleSheets.DriveFolderID,
		"credentialsFile": c.config.GoogleSheets.CredentialsFile,
	})
}

// HandleUpdateConfig menangani API untuk memperbarui konfigurasi
func (c *ConfigController) HandleUpdateConfig(ctx *fiber.Ctx) error {
	// Di implementasi sebenarnya, kita akan memvalidasi dan menyimpan perubahan
	// Namun untuk sekarang kita hanya simulasikan berhasil

	c.log.Info("Config update requested (not implemented yet)")

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Konfigurasi berhasil diperbarui",
	})
}
