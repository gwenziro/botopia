package web

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gwenziro/botopia/internal/domain/service"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
)

// ContactController adalah controller untuk manajemen kontak
type ContactController struct {
	contactService service.ContactService
	log            *logger.Logger
}

// NewContactController membuat instance contact controller baru
func NewContactController(contactService service.ContactService) *ContactController {
	return &ContactController{
		contactService: contactService,
		log:            logger.New("ContactController", logger.INFO, true),
	}
}

// HandleContactPage menangani halaman kontak
func (c *ContactController) HandleContactPage(ctx *fiber.Ctx) error {
	return ctx.Render("pages/contacts", fiber.Map{
		"Title": "Kontak & Whitelist | Botopia",
		"Page":  "contacts", // Pastikan nilainya "contacts" untuk memicu loading script
	}, "layouts/main")
}

// HandleGetContacts menangani API untuk mendapatkan daftar kontak
func (c *ContactController) HandleGetContacts(ctx *fiber.Ctx) error {
	timeoutCtx, cancel := context.WithTimeout(ctx.Context(), 5*time.Second)
	defer cancel()

	contacts, err := c.contactService.GetAllContacts(timeoutCtx)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal memuat daftar kontak: " + err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"contacts": contacts,
	})
}

// HandleGetWhitelistedContacts menangani API untuk mendapatkan kontak whitelist
func (c *ContactController) HandleGetWhitelistedContacts(ctx *fiber.Ctx) error {
	timeoutCtx, cancel := context.WithTimeout(ctx.Context(), 5*time.Second)
	defer cancel()

	contacts, err := c.contactService.GetWhitelistedContacts(timeoutCtx)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal memuat daftar whitelist: " + err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"contacts": contacts,
	})
}

// HandleAddContact menangani API untuk menambahkan kontak
func (c *ContactController) HandleAddContact(ctx *fiber.Ctx) error {
	var input struct {
		Phone    string `json:"phone"`
		Name     string `json:"name"`
		Notes    string `json:"notes"`
		IsActive bool   `json:"isActive"`
	}

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Format data tidak valid",
		})
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.Context(), 5*time.Second)
	defer cancel()

	contact, err := c.contactService.AddContact(timeoutCtx, input.Phone, input.Name, input.Notes)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Gagal menambahkan kontak: " + err.Error(),
		})
	}

	// Set whitelist status if provided
	if input.IsActive {
		err = c.contactService.SetWhitelistStatus(timeoutCtx, contact.Phone, true)
		if err != nil {
			c.log.Error("Gagal mengatur status whitelist: %v", err)
		}
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"contact": contact,
	})
}

// HandleUpdateContact menangani API untuk memperbarui kontak
func (c *ContactController) HandleUpdateContact(ctx *fiber.Ctx) error {
	var input struct {
		Phone    string `json:"phone"`
		Name     string `json:"name"`
		Notes    string `json:"notes"`
		IsActive bool   `json:"isActive"`
	}

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Format data tidak valid",
		})
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.Context(), 5*time.Second)
	defer cancel()

	contact, err := c.contactService.UpdateContact(
		timeoutCtx, input.Phone, input.Name, input.Notes, input.IsActive)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Gagal memperbarui kontak: " + err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"contact": contact,
	})
}

// HandleDeleteContact menangani API untuk menghapus kontak
func (c *ContactController) HandleDeleteContact(ctx *fiber.Ctx) error {
	var input struct {
		Phone string `json:"phone"`
	}

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Format data tidak valid",
		})
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.Context(), 5*time.Second)
	defer cancel()

	err := c.contactService.DeleteContact(timeoutCtx, input.Phone)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menghapus kontak: " + err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

// HandleSetWhitelistStatus menangani API untuk mengatur status whitelist
func (c *ContactController) HandleSetWhitelistStatus(ctx *fiber.Ctx) error {
	var input struct {
		Phone    string `json:"phone"`
		IsActive bool   `json:"isActive"`
	}

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Format data tidak valid",
		})
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.Context(), 5*time.Second)
	defer cancel()

	err := c.contactService.SetWhitelistStatus(timeoutCtx, input.Phone, input.IsActive)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengatur status whitelist: " + err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
	})
}
