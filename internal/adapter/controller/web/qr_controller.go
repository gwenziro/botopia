package web

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
	"github.com/gwenziro/botopia/internal/usecase/connection"
)

// QRController controller untuk halaman QR
type QRController struct {
	connectWhatsAppUseCase *connection.ConnectWhatsAppUseCase
	log                    *logger.Logger
	qrCode                 string
	qrMutex                sync.RWMutex
	qrChanListener         sync.Once
}

// NewQRController membuat instance QR controller
func NewQRController(connectUC *connection.ConnectWhatsAppUseCase) *QRController {
	ctrl := &QRController{
		connectWhatsAppUseCase: connectUC,
		log:                    logger.New("QRController", logger.INFO, true),
	}

	// Start QR code listener
	ctrl.startQRListener()

	return ctrl
}

// startQRListener memulai listener untuk QR code
func (c *QRController) startQRListener() {
	c.qrChanListener.Do(func() {
		go func() {
			c.log.Info("Starting QR code listener")
			qrChan := c.connectWhatsAppUseCase.GetQRChannel()

			for qrCode := range qrChan {
				if qrCode == "" {
					continue
				}
				c.UpdateQRCode(qrCode)
			}
		}()
	})
}

// HandleQRPage menangani halaman QR
func (c *QRController) HandleQRPage(ctx *fiber.Ctx) error {
	// Get current QR code
	c.qrMutex.RLock()
	qrCode := c.qrCode
	c.qrMutex.RUnlock()

	// Get connection status
	isConnected := c.connectWhatsAppUseCase.IsConnected()
	connState := "disconnected"
	if isConnected {
		connState = "connected"
	}

	// Get phone number if connected
	var phone string
	if isConnected {
		// In a real implementation, we would get this from the use case
		// For now, we'll leave it empty
		phone = ""
	}

	return ctx.Render("pages/qr", fiber.Map{
		"Title":           "QR Code | Botopia",
		"Page":            "qr",
		"QRCode":          qrCode,
		"ConnectionState": connState,
		"Phone":           phone,
		"Name":            "WhatsApp User",
	}, "layouts/main")
}

// HandleGetQR menangani API untuk QR code
func (c *QRController) HandleGetQR(ctx *fiber.Ctx) error {
	// Get current QR code
	c.qrMutex.RLock()
	qrCode := c.qrCode
	c.qrMutex.RUnlock()

	// Get connection status
	isConnected := c.connectWhatsAppUseCase.IsConnected()

	// Inisialisasi response
	response := fiber.Map{
		"qrCode":          qrCode,
		"connectionState": isConnected,
	}

	// Get phone number and device info if connected
	if isConnected {
		// Get detailed user info
		userInfo, err := c.connectWhatsAppUseCase.GetCurrentUser()
		if err != nil {
			c.log.Warn("Error getting user info: %v", err)
		}

		if userInfo != nil {
			response["phone"] = userInfo.Phone
			response["name"] = userInfo.Name
			if userInfo.Name == "" {
				response["name"] = "WhatsApp User"
			}

			// Add device details if available
			if userInfo.DeviceDetails != nil {
				response["deviceInfo"] = userInfo.DeviceDetails
				response["platform"] = userInfo.DeviceDetails.Platform
				response["deviceID"] = userInfo.DeviceDetails.DeviceID
				response["connectedSince"] = userInfo.DeviceDetails.Connected
				response["clientIP"] = ctx.IP()

				if userInfo.DeviceDetails.BusinessName != "" {
					response["businessName"] = userInfo.DeviceDetails.BusinessName
				}
			}
		} else {
			response["phone"] = ""
			response["name"] = "WhatsApp User"
			response["clientIP"] = ctx.IP()
		}
	}

	return ctx.JSON(response)
}

// UpdateQRCode memperbarui kode QR tersimpan
func (c *QRController) UpdateQRCode(code string) {
	c.qrMutex.Lock()
	defer c.qrMutex.Unlock()

	if code == "" {
		c.log.Warn("Tried to save empty QR code, ignoring")
		return
	}

	// Don't update if already connected
	if c.connectWhatsAppUseCase.IsConnected() {
		c.log.Info("Already connected, ignoring QR code update")
		return
	}

	c.qrCode = code
	c.log.Info("Updated QR code (length: %d characters)", len(code))
}

// HandleDisconnect menangani API untuk memutuskan koneksi
func (c *QRController) HandleDisconnect(ctx *fiber.Ctx) error {
	c.log.Info("Memutuskan koneksi WhatsApp via API")

	err := c.connectWhatsAppUseCase.Disconnect()
	if err != nil {
		c.log.Error("Gagal memutuskan koneksi: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Reset QR code
	c.qrMutex.Lock()
	c.qrCode = ""
	c.qrMutex.Unlock()

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Koneksi berhasil diputuskan",
	})
}
