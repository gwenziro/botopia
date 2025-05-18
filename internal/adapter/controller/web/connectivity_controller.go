package web

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gwenziro/botopia/internal/usecase/connection"
	"github.com/gwenziro/botopia/internal/usecase/stats"
)

// ConnectivityController adalah controller untuk halaman connectivity dan koneksi WhatsApp
type ConnectivityController struct {
	connectUseCase *connection.ConnectWhatsAppUseCase
	statsUseCase   *stats.GetStatsUseCase
}

// NewConnectivityController membuat instance connectivity controller baru
func NewConnectivityController(
	connectUC *connection.ConnectWhatsAppUseCase,
	statsUC *stats.GetStatsUseCase,
) *ConnectivityController {
	return &ConnectivityController{
		connectUseCase: connectUC,
		statsUseCase:   statsUC,
	}
}

// HandleConnectivityPage menangani halaman connectivity
func (c *ConnectivityController) HandleConnectivityPage(ctx *fiber.Ctx) error {
	return ctx.Render("pages/qr", fiber.Map{
		"Title": "Connectivity | Botopia",
		"Page":  "connectivity",
	}, "layouts/main")
}

// HandleGetQR menangani API untuk mendapatkan QR code
func (c *ConnectivityController) HandleGetQR(ctx *fiber.Ctx) error {
	// Check if already connected
	if c.connectUseCase.IsConnected() {
		user, err := c.connectUseCase.GetCurrentUser()
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get user: " + err.Error(),
			})
		}

		phone := ""
		name := "WhatsApp User"

		if user != nil {
			phone = user.Phone
			if user.Name != "" {
				name = user.Name
			} else if user.PushName != "" {
				name = user.PushName
			}
		}

		return ctx.JSON(fiber.Map{
			"qrCode":          "",
			"connectionState": true,
			"phone":           phone,
			"name":            name,
		})
	}

	// Create a channel for QR code
	qrChan := c.connectUseCase.GetQRChannel()

	// Wait for QR code with timeout
	select {
	case <-ctx.Context().Done():
		return ctx.Status(http.StatusRequestTimeout).JSON(fiber.Map{
			"error": "QR code generation timed out",
		})
	case qrCode := <-qrChan:
		if qrCode == "" {
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get QR code",
			})
		}

		return ctx.JSON(fiber.Map{
			"qrCode":          qrCode,
			"connectionState": false,
		})
	}
}

// HandleDisconnect menangani pemutus koneksi WhatsApp
func (c *ConnectivityController) HandleDisconnect(ctx *fiber.Ctx) error {
	// Cek apakah terhubung
	if !c.connectUseCase.IsConnected() {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Tidak terhubung ke WhatsApp",
		})
	}

	// Disconnect
	err := c.connectUseCase.Disconnect()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal memutuskan koneksi: " + err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Koneksi WhatsApp berhasil diputus",
	})
}
