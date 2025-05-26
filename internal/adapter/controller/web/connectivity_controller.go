package web

import (
	"context"
	"log"
	"net/http"
	"time"

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
	return ctx.Render("pages/connectivity", fiber.Map{
		"Title": "Konektivitas | Botopia",
		"Page":  "connectivity",
	}, "layouts/main")
}

// HandleGetQR menangani API untuk mendapatkan QR code
func (c *ConnectivityController) HandleGetQR(ctx *fiber.Ctx) error {
	// Check if already connected
	if c.connectUseCase.IsConnected() {
		// Get connection stats from stats use case
		stats, err := c.statsUseCase.Execute(ctx.Context())
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get stats: " + err.Error(),
			})
		}

		// Get current user info
		user, err := c.connectUseCase.GetCurrentUser()
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get user: " + err.Error(),
			})
		}

		phone := ""
		name := "WhatsApp User"
		var deviceDetails map[string]interface{} = nil

		if user != nil {
			phone = user.Phone
			if user.Name != "" {
				name = user.Name
			} else if user.PushName != "" {
				name = user.PushName
			}

			// Provide simplified device details
			if user.DeviceDetails != nil {
				deviceDetails = map[string]interface{}{
					"platform":  user.DeviceDetails.Platform,
					"deviceId":  user.DeviceDetails.DeviceID,
					"ipAddress": user.DeviceDetails.IPAddress,
				}
			}
		}

		return ctx.JSON(fiber.Map{
			"qrCode":          "",
			"connectionState": true,
			"phone":           phone,
			"name":            name,
			"uptime":          stats.Uptime,
			"deviceDetails":   deviceDetails,
		})
	}

	// Mulai koneksi jika belum terhubung dan belum mencoba terhubung
	if !c.connectUseCase.IsConnecting() {
		go func() {
			_, err := c.connectUseCase.Execute(context.Background())
			if err != nil {
				log.Printf("Error starting connection: %v", err)
			}
		}()
	}

	// Dapatkan QR code yang sudah tersedia (jika ada)
	qrCode := c.connectUseCase.GetCurrentQR()
	if qrCode != "" {
		return ctx.JSON(fiber.Map{
			"qrCode":          qrCode,
			"connectionState": false,
		})
	}

	// Create a channel for QR code dengan buffer
	qrChan := c.connectUseCase.GetQRChannel()

	// Kurangi timeout dari 15 detik menjadi 5 detik
	select {
	case <-ctx.Context().Done():
		return ctx.Status(http.StatusRequestTimeout).JSON(fiber.Map{
			"error":           "QR code generation timed out",
			"connectionState": false,
			"retry":           true,
		})
	case qrCode := <-qrChan:
		if qrCode == "" {
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error":           "QR code unavailable",
				"connectionState": false,
				"retry":           true,
			})
		}

		return ctx.JSON(fiber.Map{
			"qrCode":          qrCode,
			"connectionState": false,
		})
	case <-time.After(5 * time.Second):
		return ctx.JSON(fiber.Map{
			"connectionState": false,
			"retry":           true,
			"message":         "QR code not yet available, please retry",
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
