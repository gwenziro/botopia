package web

import (
	"fmt"

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
	return ctx.Render("pages/config", fiber.Map{
		"Title":             "Konfigurasi | Botopia",
		"Page":              "konfigurasi",
		"LogLevel":          c.config.LogLevel,
		"CommandPrefix":     c.config.CommandPrefix,
		"WebPort":           c.config.WebPort,
		"WebHost":           c.config.WebHost,
		"WebAuthEnabled":    c.config.WebAuthEnabled,
		"WebAuthUsername":   c.config.WebAuthUsername,
		"WebAuthPassword":   c.config.WebAuthPassword,
		"DevMode":           c.config.DevMode,
		"QRAutoRefresh":     c.config.QRAutoRefresh,
		"GoogleCredentials": c.config.GoogleSheets.CredentialsFile,
		"SpreadsheetID":     c.config.GoogleSheets.SpreadsheetID,
		"DriveFolderID":     c.config.GoogleSheets.DriveFolderID,
	}, "layouts/main")
}

// HandleGetConfig menangani API untuk mendapatkan konfigurasi
func (c *ConfigController) HandleGetConfig(ctx *fiber.Ctx) error {
	// Sembunyikan password sebelum mengembalikan konfigurasi
	config := map[string]interface{}{
		"logLevel":        c.config.LogLevel,
		"commandPrefix":   c.config.CommandPrefix,
		"webPort":         c.config.WebPort,
		"webHost":         c.config.WebHost,
		"webAuthEnabled":  c.config.WebAuthEnabled,
		"webAuthUsername": c.config.WebAuthUsername,
		"webAuthPassword": "********", // Sembunyikan password
		"devMode":         c.config.DevMode,
		"qrAutoRefresh":   c.config.QRAutoRefresh,
		"googleSheets": map[string]string{
			"credentialsFile": c.config.GoogleSheets.CredentialsFile,
			"spreadsheetID":   c.config.GoogleSheets.SpreadsheetID,
			"driveFolderID":   c.config.GoogleSheets.DriveFolderID,
		},
	}

	return ctx.JSON(config)
}

// HandleUpdateConfig menangani API untuk memperbarui konfigurasi
func (c *ConfigController) HandleUpdateConfig(ctx *fiber.Ctx) error {
	var configInput struct {
		CommandPrefix   string `json:"commandPrefix"`
		LogLevel        string `json:"logLevel"`
		WebPort         int    `json:"webPort"`
		WebHost         string `json:"webHost"`
		WebAuthEnabled  bool   `json:"webAuthEnabled"`
		WebAuthUsername string `json:"webAuthUsername"`
		WebAuthPassword string `json:"webAuthPassword,omitempty"`
		GoogleSheets    struct {
			SpreadsheetID   string `json:"spreadsheetID"`
			DriveFolderID   string `json:"driveFolderID"`
			CredentialsFile string `json:"credentialsFile"`
		} `json:"googleSheets"`
	}

	if err := ctx.BodyParser(&configInput); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request format",
		})
	}

	// Log received data for debugging
	c.log.Info("Received config update: %+v", configInput)

	// Validasi command prefix hanya terima ! atau #
	if configInput.CommandPrefix != "!" && configInput.CommandPrefix != "#" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Command prefix hanya bisa '!' atau '#'",
		})
	}

	// Update config values
	c.config.CommandPrefix = configInput.CommandPrefix
	c.config.LogLevel = configInput.LogLevel
	c.config.WebPort = configInput.WebPort
	c.config.WebHost = configInput.WebHost
	c.config.WebAuthEnabled = configInput.WebAuthEnabled
	c.config.WebAuthUsername = configInput.WebAuthUsername

	if configInput.WebAuthPassword != "" && configInput.WebAuthPassword != "********" {
		c.config.WebAuthPassword = configInput.WebAuthPassword
	}

	// Update Google Sheets Config
	c.config.GoogleSheets.SpreadsheetID = configInput.GoogleSheets.SpreadsheetID
	c.config.GoogleSheets.DriveFolderID = configInput.GoogleSheets.DriveFolderID
	c.config.GoogleSheets.CredentialsFile = configInput.GoogleSheets.CredentialsFile

	// TODO: Simpan perubahan konfigurasi ke file atau environment

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Configuration updated successfully",
	})
}

// HandleGetConfigStatus returns configuration status as JSON
func (c *ConfigController) HandleGetConfigStatus(ctx *fiber.Ctx) error {
	// Check Google API status
	isGoogleAPIConfigured := c.config.GoogleSheets.CredentialsFile != ""
	googleSheetsConfigured := c.config.GoogleSheets.SpreadsheetID != ""
	googleDriveConfigured := c.config.GoogleSheets.DriveFolderID != ""

	// Get spreadsheet URL & ID
	spreadsheetUrl := ""
	spreadsheetId := c.config.GoogleSheets.SpreadsheetID
	if googleSheetsConfigured {
		spreadsheetUrl = fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/edit", spreadsheetId)
	}

	// Get drive folder URL & ID
	driveFolderUrl := ""
	driveFolderId := c.config.GoogleSheets.DriveFolderID
	if googleDriveConfigured {
		driveFolderUrl = fmt.Sprintf("https://drive.google.com/drive/folders/%s", driveFolderId)
	}

	return ctx.JSON(fiber.Map{
		"googleApi": fiber.Map{
			"configured": isGoogleAPIConfigured,
			"sheets":     googleSheetsConfigured,
			"drive":      googleDriveConfigured,
		},
		"spreadsheetUrl": spreadsheetUrl,
		"spreadsheetId":  spreadsheetId,
		"driveFolderUrl": driveFolderUrl,
		"driveFolderId":  driveFolderId,
		"systemStatus": fiber.Map{
			"devMode":       c.config.DevMode,
			"qrAutoRefresh": c.config.QRAutoRefresh,
			"version":       "1.0.0", // Tambahkan versi statis untuk informasi sementara
		},
	})
}
