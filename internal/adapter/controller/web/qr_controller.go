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

	// Get phone number if connected
	var phone string
	if isConnected {
		// In a real implementation, we would get this from the use case
		phone = ""
	}

	return ctx.JSON(fiber.Map{
		"qrCode":          qrCode,
		"connectionState": isConnected,
		"phone":           phone,
		"name":            "WhatsApp User",
	})
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
