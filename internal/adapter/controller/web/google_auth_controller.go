package web

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gwenziro/botopia/internal/domain/repository"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
)

// GoogleAuthController adalah controller untuk otorisasi Google API
type GoogleAuthController struct {
	log     *logger.Logger
	apiRepo repository.GoogleAPIRepository
}

// NewGoogleAuthController membuat instance controller baru
func NewGoogleAuthController(apiRepo repository.GoogleAPIRepository) *GoogleAuthController {
	return &GoogleAuthController{
		log:     logger.New("GoogleAuthController", logger.INFO, true),
		apiRepo: apiRepo,
	}
}

// HandleAuthStart menangani permintaan untuk memulai otorisasi
func (c *GoogleAuthController) HandleAuthStart(ctx *fiber.Ctx) error {
	c.log.Info("Auth start requested, menggunakan Service Account")

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message":         "Aplikasi menggunakan Service Account untuk otentikasi Google API. Anda tidak perlu melakukan otorisasi.",
		"status":          "service_account_active",
		"configured":      c.apiRepo.IsConfigured(),
		"credentialsPath": c.apiRepo.GetCredentialsPath(),
	})
}

// HandleAuthCallback menangani callback dari Google setelah otorisasi
func (c *GoogleAuthController) HandleAuthCallback(ctx *fiber.Ctx) error {
	c.log.Info("Auth callback dipanggil, menggunakan Service Account")

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Aplikasi menggunakan Service Account untuk otentikasi Google API. Anda tidak perlu melakukan otorisasi.",
		"status":  "service_account_active",
	})
}
