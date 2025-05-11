package web

import (
	"github.com/gofiber/fiber/v2"
)

// AuthController controller untuk autentikasi
type AuthController struct {
	username    string
	password    string
	authEnabled bool
}

// NewAuthController membuat instance auth controller
func NewAuthController(username, password string, authEnabled bool) *AuthController {
	return &AuthController{
		username:    username,
		password:    password,
		authEnabled: authEnabled,
	}
}

// HandleLogin menangani halaman login
func (c *AuthController) HandleLogin(ctx *fiber.Ctx) error {
	redirect := ctx.Query("redirect", "/dashboard")

	return ctx.Render("pages/login", fiber.Map{
		"Title":    "Login | Botopia",
		"Page":     "login",
		"Redirect": redirect,
		"Error":    ctx.Query("error"),
	}, "layouts/minimal")
}

// HandleLoginPost menangani proses login
func (c *AuthController) HandleLoginPost(ctx *fiber.Ctx) error {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")
	redirect := ctx.FormValue("redirect", "/dashboard")

	if username == c.username && password == c.password {
		// Set cookie autentikasi
		ctx.Cookie(&fiber.Cookie{
			Name:     "authenticated",
			Value:    "true",
			Path:     "/",
			MaxAge:   86400, // 24 jam
			HTTPOnly: true,
			Secure:   ctx.Secure(),
			SameSite: "Lax",
		})

		return ctx.Redirect(redirect)
	}

	// Login gagal
	return ctx.Redirect("/login?error=Invalid+username+or+password&redirect=" + redirect)
}

// IsAuthEnabled memeriksa apakah autentikasi diaktifkan
func (c *AuthController) IsAuthEnabled() bool {
	return c.authEnabled
}

// GetUsername mendapatkan username yang dikonfigurasi
func (c *AuthController) GetUsername() string {
	return c.username
}

// GetPassword mendapatkan password yang dikonfigurasi
func (c *AuthController) GetPassword() string {
	return c.password
}
