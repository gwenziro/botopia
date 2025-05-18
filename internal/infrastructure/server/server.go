package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
	configCtrl := s.container.GetConfigController()
	dataMasterCtrl := s.container.GetDataMasterController()

	// Setup message controller untuk WhatsApp
	s.container.GetMessageController().Setup()

	// Setup routes based on auth status
	if authCtrl.IsAuthEnabled() {
		s.setupAuthenticatedRoutes(dashboardCtrl, qrCtrl, authCtrl, configCtrl)
	} else {
		s.setupUnauthenticatedRoutes(dashboardCtrl, qrCtrl, configCtrl, dataMasterCtrl, s.container.GetContactController())
	}

	s.log.Info("Routes set up successfully")
}

// Start memulai web server
func (s *Server) Start() error {
	// Setup routes
	s.SetupRoutes()

	addr := fmt.Sprintf("%s:%d", s.container.GetConfig().WebHost, s.port)
	s.log.Info("Starting web server on %s", addr)

	// Jalankan server secara non-blocking
	return s.app.Listen(addr)
}

// Shutdown menghentikan web server secara graceful
func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down web server")
	return s.app.ShutdownWithContext(ctx)
}

// setupAuthenticatedRoutes mengatur route dengan auth
func (s *Server) setupAuthenticatedRoutes(
	dashboardCtrl, qrCtrl, authCtrl, configCtrl interface{},
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

	config := configCtrl.(interface {
		HandleConfigPage(ctx *fiber.Ctx) error
		HandleGetConfig(ctx *fiber.Ctx) error
		HandleUpdateConfig(ctx *fiber.Ctx) error
		HandleGetConfigStatus(ctx *fiber.Ctx) error
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
	s.app.Get("/konfigurasi", authMiddleware, config.HandleConfigPage)

	// API routes
	api := s.app.Group("/api", authMiddleware)
	api.Get("/stats", dashboard.HandleGetStats)
	api.Get("/qr", qr.HandleGetQR)
	api.Get("/config", config.HandleGetConfig)
	api.Post("/config", config.HandleUpdateConfig)
	api.Get("/config/status", config.HandleGetConfigStatus)
}

// setupUnauthenticatedRoutes mengatur route tanpa auth
func (s *Server) setupUnauthenticatedRoutes(
	dashboardCtrl, qrCtrl, configCtrl, dataMasterCtrl, contactCtrl interface{},
) {
	// Cast to correct types
	dashboard := dashboardCtrl.(interface {
		HandleIndex(ctx *fiber.Ctx) error
		HandleDashboard(ctx *fiber.Ctx) error
		HandleGetStats(ctx *fiber.Ctx) error
		HandleGetCommands(ctx *fiber.Ctx) error
		HandleGetRecentTransactions(ctx *fiber.Ctx) error // Add this
	})

	qr := qrCtrl.(interface {
		HandleQRPage(ctx *fiber.Ctx) error
		HandleGetQR(ctx *fiber.Ctx) error
		HandleDisconnect(ctx *fiber.Ctx) error
	})

	config := configCtrl.(interface {
		HandleConfigPage(ctx *fiber.Ctx) error
		HandleGetConfig(ctx *fiber.Ctx) error
		HandleUpdateConfig(ctx *fiber.Ctx) error
		HandleGetConfigStatus(ctx *fiber.Ctx) error // Add this to the interface
	})

	dataMaster := dataMasterCtrl.(interface {
		HandleDataMasterPage(ctx *fiber.Ctx) error
		HandleGetMasterData(ctx *fiber.Ctx) error
	})

	contact := contactCtrl.(interface {
		HandleContactPage(ctx *fiber.Ctx) error
		HandleGetContacts(ctx *fiber.Ctx) error
		HandleGetWhitelistedContacts(ctx *fiber.Ctx) error
		HandleAddContact(ctx *fiber.Ctx) error
		HandleUpdateContact(ctx *fiber.Ctx) error
		HandleDeleteContact(ctx *fiber.Ctx) error
		HandleSetWhitelistStatus(ctx *fiber.Ctx) error
	})

	// Public routes
	s.app.Get("/", dashboard.HandleIndex)
	s.app.Get("/dashboard", dashboard.HandleDashboard)
	s.app.Get("/qr", qr.HandleQRPage)
	s.app.Get("/konfigurasi", config.HandleConfigPage)
	s.app.Get("/data-master", dataMaster.HandleDataMasterPage) // Tambahkan route data master
	s.app.Get("/contacts", contact.HandleContactPage)

	// Cast Commands controller
	commands := s.container.GetCommandsController()

	// Tambahkan route baru untuk commands
	s.app.Get("/commands", commands.HandleCommandsPage)

	// API routes
	api := s.app.Group("/api")
	api.Get("/stats", dashboard.HandleGetStats)
	api.Get("/commands", dashboard.HandleGetCommands)
	api.Get("/transactions/recent", dashboard.HandleGetRecentTransactions) // Add this
	api.Get("/qr", qr.HandleGetQR)
	api.Post("/disconnect", qr.HandleDisconnect)
	api.Get("/config", config.HandleGetConfig)
	api.Post("/config", config.HandleUpdateConfig)
	api.Get("/config/status", config.HandleGetConfigStatus) // Add this new route

	// Data Master API routes - hanya route GET yang diperlukan
	api.Get("/data-master", dataMaster.HandleGetMasterData)

	// Contact API routes
	api.Get("/contacts", contact.HandleGetContacts)
	api.Get("/contacts/whitelist", contact.HandleGetWhitelistedContacts)
	api.Post("/contacts/add", contact.HandleAddContact)
	api.Post("/contacts/update", contact.HandleUpdateContact)
	api.Post("/contacts/delete", contact.HandleDeleteContact)
	api.Post("/whitelist/status", contact.HandleSetWhitelistStatus)

	// Whitelist Toggle API
	api.Get("/whitelist/status", func(ctx *fiber.Ctx) error {
		// Dapatkan status whitelist dari MessageController
		return ctx.JSON(fiber.Map{
			"enabled": s.container.GetMessageController().UseWhitelist,
		})
	})

	api.Post("/whitelist/toggle", func(ctx *fiber.Ctx) error {
		var input struct {
			Enabled bool `json:"enabled"`
		}

		if err := ctx.BodyParser(&input); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Format data tidak valid",
			})
		}

		// Update status whitelist di message controller
		s.container.GetMessageController().SetUseWhitelist(input.Enabled)

		return ctx.JSON(fiber.Map{
			"success": true,
			"enabled": input.Enabled,
		})
	})
}
