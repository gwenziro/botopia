package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware adalah middleware untuk autentikasi
type AuthMiddleware struct {
	// Username dan password untuk autentikasi
	Username string
	Password string
}

// NewAuthMiddleware membuat instance AuthMiddleware baru
func NewAuthMiddleware(username, password string) *AuthMiddleware {
	return &AuthMiddleware{
		Username: username,
		Password: password,
	}
}

// RequireAuth adalah middleware untuk mengharuskan autentikasi
func (m *AuthMiddleware) RequireAuth(c *fiber.Ctx) error {
	// Check if authenticated via cookie
	if c.Cookies("authenticated") == "true" {
		return c.Next()
	}

	// Redirect to login
	return c.Redirect("/login?redirect=" + c.Path())
}

// BasicAuth adalah middleware untuk basic auth
func (m *AuthMiddleware) BasicAuth(c *fiber.Ctx) error {
	// Implement basic auth here
	// ...

	return c.Next()
}

// Implementasi ini masih diperlukan untuk menggantikan middleware yang dibuat langsung di server.go
