package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"

	"github.com/gwenziro/botopia/internal/app/di"
	botLogger "github.com/gwenziro/botopia/internal/infrastructure/logger"
)

// Server mengelola server web
type Server struct {
	app       *fiber.App
	container *di.Container
	log       *botLogger.Logger
	port      int
}

// NewServer membuat instance server baru
func NewServer(container *di.Container) *Server {
	log := botLogger.New("WebServer", botLogger.INFO, true)

	// Setup template engine - pastikan direktori view sudah benar
	viewDir := "./internal/infrastructure/web/view"
	staticDir := "./internal/infrastructure/web/static"

	// Ini memberi kesempatan untuk migrasi bertahap
	engine := html.New(viewDir, ".html")

	// Add template functions
	engine.AddFunc("json", func(v interface{}) (string, error) {
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(b), nil
	})

	engine.AddFunc("dict", func(values ...interface{}) (map[string]interface{}, error) {
		if len(values)%2 != 0 {
			return nil, fmt.Errorf("dict requires an even number of arguments")
		}
		dict := make(map[string]interface{}, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			key, ok := values[i].(string)
			if !ok {
				return nil, fmt.Errorf("dict keys must be strings")
			}
			dict[key] = values[i+1]
		}
		return dict, nil
	})

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New())

	// Static files - cek direktori dan fallback ke direktori lama jika perlu
	app.Static("/static", staticDir)

	return &Server{
		app:       app,
		container: container,
		log:       log,
		port:      container.GetConfig().WebPort,
	}
}

// SetupRoutes mengatur route aplikasi
func (s *Server) SetupRoutes() {
	// Get controllers from container
	dashboardCtrl := s.container.GetDashboardController()
	qrCtrl := s.container.GetQRController()
	authCtrl := s.container.GetAuthController()

	// Setup message controller untuk WhatsApp
	s.container.GetMessageController().Setup()

	// Setup routes based on auth status
	if authCtrl.IsAuthEnabled() {
		s.setupAuthenticatedRoutes(dashboardCtrl, qrCtrl, authCtrl)
	} else {
		s.setupUnauthenticatedRoutes(dashboardCtrl, qrCtrl)
	}

	s.log.Info("Routes set up successfully")
}

// setupAuthenticatedRoutes mengatur route dengan auth
func (s *Server) setupAuthenticatedRoutes(
	dashboardCtrl, qrCtrl, authCtrl interface{},
) {
	// Cast to correct types
	dashboard := dashboardCtrl.(interface {
		HandleIndex(ctx *fiber.Ctx) error
		HandleDashboard(ctx *fiber.Ctx) error
		HandleGetStats(ctx *fiber.Ctx) error
	})

	qr := qrCtrl.(interface {
		HandleQRPage(ctx *fiber.Ctx) error
		HandleGetQR(ctx *fiber.Ctx) error
	})

	auth := authCtrl.(interface {
		HandleLogin(ctx *fiber.Ctx) error
		HandleLoginPost(ctx *fiber.Ctx) error
	})

	// Auth routes
	s.app.Get("/login", auth.HandleLogin)
	s.app.Post("/login", auth.HandleLoginPost)

	// Auth middleware
	authMiddleware := func(c *fiber.Ctx) error {
		if c.Cookies("authenticated") != "true" {
			return c.Redirect("/login?redirect=" + c.Path())
		}
		return c.Next()
	}

	// Protected routes
	s.app.Get("/", authMiddleware, dashboard.HandleIndex)
	s.app.Get("/dashboard", authMiddleware, dashboard.HandleDashboard)
	s.app.Get("/qr", authMiddleware, qr.HandleQRPage)

	// API routes
	api := s.app.Group("/api", authMiddleware)
	api.Get("/stats", dashboard.HandleGetStats)
	api.Get("/qr", qr.HandleGetQR)
}

// setupUnauthenticatedRoutes mengatur route tanpa auth
func (s *Server) setupUnauthenticatedRoutes(
	dashboardCtrl, qrCtrl interface{},
) {
	// Cast to correct types
	dashboard := dashboardCtrl.(interface {
		HandleIndex(ctx *fiber.Ctx) error
		HandleDashboard(ctx *fiber.Ctx) error
		HandleGetStats(ctx *fiber.Ctx) error
	})

	qr := qrCtrl.(interface {
		HandleQRPage(ctx *fiber.Ctx) error
		HandleGetQR(ctx *fiber.Ctx) error
	})

	// Public routes
	s.app.Get("/", dashboard.HandleIndex)
	s.app.Get("/dashboard", dashboard.HandleDashboard)
	s.app.Get("/qr", qr.HandleQRPage)

	// API routes
	api := s.app.Group("/api")
	api.Get("/stats", dashboard.HandleGetStats)
	api.Get("/qr", qr.HandleGetQR)
}

// Start memulai server web
func (s *Server) Start() error {
	// Setup routes
	s.SetupRoutes()

	// Start server
	s.log.Info("Starting web server on port %d", s.port)
	return s.app.Listen(fmt.Sprintf(":%d", s.port))
}

// Shutdown menghentikan server
func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down web server")
	return s.app.ShutdownWithContext(ctx)
}
