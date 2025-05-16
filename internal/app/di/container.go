package di

import (
	"database/sql"
	"fmt"

	"github.com/gwenziro/botopia/internal/adapter/controller/web"
	whatsappController "github.com/gwenziro/botopia/internal/adapter/controller/whatsapp"
	googleRepo "github.com/gwenziro/botopia/internal/adapter/repository/google"
	"github.com/gwenziro/botopia/internal/adapter/repository/memory"
	whatsmeowRepo "github.com/gwenziro/botopia/internal/adapter/repository/whatsmeow"
	adapterService "github.com/gwenziro/botopia/internal/adapter/service"
	"github.com/gwenziro/botopia/internal/app/command"
	"github.com/gwenziro/botopia/internal/domain/repository"
	"github.com/gwenziro/botopia/internal/domain/service"
	"github.com/gwenziro/botopia/internal/infrastructure/config"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
	"github.com/gwenziro/botopia/internal/usecase/command/execute"
	"github.com/gwenziro/botopia/internal/usecase/command/list"
	connectionUseCase "github.com/gwenziro/botopia/internal/usecase/connection"
	"github.com/gwenziro/botopia/internal/usecase/stats"

	// Dibutuhkan untuk SQLite
	_ "modernc.org/sqlite"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// Container adalah container untuk dependency injection
type Container struct {
	config *config.Config
	log    *logger.Logger
	db     *sql.DB

	// Repositories
	commandRepository    *memory.CommandRepository
	connectionRepository *whatsmeowRepo.ConnectionRepository
	statsRepository      *memory.StatsRepository
	googleAPIRepository  *googleRepo.GoogleAPIRepository
	sheetsRepository     *googleRepo.SheetsRepository
	driveRepository      *googleRepo.DriveRepository
	contactRepository    *memory.ContactRepository

	// Use cases
	executeCommandUseCase  *execute.ExecuteCommandUseCase
	listCommandsUseCase    *list.ListCommandsUseCase
	connectWhatsAppUseCase *connectionUseCase.ConnectWhatsAppUseCase
	getStatsUseCase        *stats.GetStatsUseCase

	// Services
	financeService service.FinanceService
	contactService service.ContactService

	// Controllers
	dashboardController  *web.DashboardController
	qrController         *web.QRController
	authController       *web.AuthController
	messageController    *whatsappController.MessageController
	configController     *web.ConfigController
	dataMasterController *web.DataMasterController
	contactController    *web.ContactController

	// Command initializer
	commandInitializer *command.CommandInitializer
}

// NewContainer membuat container baru
func NewContainer(cfg *config.Config) *Container {
	c := &Container{
		config: cfg,
	}

	// Pastikan direktori yang diperlukan ada
	if err := cfg.EnsureDirectories(); err != nil {
		panic(fmt.Sprintf("Failed to ensure directories: %v", err))
	}

	// Initialize components
	c.initLogger()
	c.initDatabase()
	c.initRepositories()
	c.initServices()           // Menggunakan initServices sebagai pengganti initMediaServices
	c.initCommandInitializer() // Tambahkan fungsi yang hilang
	c.initUseCases()
	c.initControllers()

	return c
}

// initLogger menginisialisasi logger
func (c *Container) initLogger() {
	c.log = logger.New("App", logger.LevelFromString(c.config.LogLevel), c.config.UseColors)
}

// initDatabase menginisialisasi koneksi database
func (c *Container) initDatabase() {
	var err error

	dsn := fmt.Sprintf("file:%s?_journal_mode=WAL&_busy_timeout=5000", c.config.DBPath)
	c.db, err = sql.Open("sqlite", dsn)
	if err != nil {
		c.log.Fatal("Gagal membuka database: %v", err)
	}

	// Pengaturan koneksi pool
	c.db.SetMaxOpenConns(10)
	c.db.SetMaxIdleConns(5)
}

// initRepositories menginisialisasi repositories
func (c *Container) initRepositories() {
	// Buat WhatsApp client
	waLogger := &whatsmeowLoggerAdapter{c.log}
	container := sqlstore.NewWithDB(c.db, "sqlite", waLogger)
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		c.log.Warn("Tidak menemukan device: %v", err)
		// Lanjutkan meskipun tidak ada device, bisa dibuat nanti dengan QR
		deviceStore = container.NewDevice()
	}
	client := whatsmeow.NewClient(deviceStore, waLogger)

	// Inisialisasi repositories
	c.commandRepository = memory.NewCommandRepository()
	c.statsRepository = memory.NewStatsRepository()
	c.connectionRepository = whatsmeowRepo.NewConnectionRepository(client, c.log, c.statsRepository)

	// Inisialisasi Google API Repository
	c.googleAPIRepository = googleRepo.NewGoogleAPIRepository(c.config, c.log)

	// Inisialisasi Google Sheets Repository
	c.sheetsRepository = googleRepo.NewSheetsRepository(c.googleAPIRepository, c.config, c.log)

	// Inisialisasi Google Drive Repository
	c.driveRepository = googleRepo.NewDriveRepository(c.googleAPIRepository, c.config, c.log)

	// Tambahkan contact repository
	c.contactRepository = memory.NewContactRepository(c.log)
}

// initServices menginisialisasi layanan
func (c *Container) initServices() {
	// Inisialisasi finance service
	c.financeService = adapterService.NewFinanceService(
		c.sheetsRepository,
		c.driveRepository,
		c.log,
	)

	// Inisialisasi contact service
	c.contactService = adapterService.NewContactService(
		c.contactRepository,
		c.log,
	)

	c.log.Info("Services berhasil diinisialisasi")
}

// initCommandInitializer menginisialisasi command initializer
func (c *Container) initCommandInitializer() {
	c.commandInitializer = command.NewCommandInitializer(c.commandRepository, c.financeService)
	c.commandInitializer.RegisterDefaultCommands()
	c.log.Info("Command default berhasil didaftarkan. Total: %d command",
		c.commandInitializer.GetCommandCount())
}

// initUseCases menginisialisasi use cases
func (c *Container) initUseCases() {
	c.executeCommandUseCase = execute.NewExecuteCommandUseCase(
		c.commandRepository,
		c.statsRepository,
		c.connectionRepository,
		c.config.CommandPrefix,
	)

	c.listCommandsUseCase = list.NewListCommandsUseCase(c.commandRepository)

	c.connectWhatsAppUseCase = connectionUseCase.NewConnectWhatsAppUseCase(c.connectionRepository)

	c.getStatsUseCase = stats.NewGetStatsUseCase(
		c.statsRepository,
		c.connectionRepository,
		c.commandRepository,
	)
}

// initControllers menginisialisasi controllers
func (c *Container) initControllers() {
	c.dashboardController = web.NewDashboardController(
		c.getStatsUseCase,
		c.listCommandsUseCase,
	)

	c.qrController = web.NewQRController(c.connectWhatsAppUseCase)

	c.authController = web.NewAuthController(
		c.config.WebAuthUsername,
		c.config.WebAuthPassword,
		c.config.WebAuthEnabled,
	)

	// Tambahkan controller kontak
	c.contactController = web.NewContactController(c.contactService)

	// Perbarui message controller dengan contactService
	c.messageController = whatsappController.NewMessageController(
		c.executeCommandUseCase,
		c.connectionRepository,
		c.statsRepository,
		c.contactService,
	)

	c.configController = web.NewConfigController(c.config)

	// Data Master controller
	c.dataMasterController = web.NewDataMasterController(c.financeService)

	c.log.Info("Controllers berhasil diinisialisasi")
}

// GetDB mengembalikan koneksi database
func (c *Container) GetDB() *sql.DB {
	return c.db
}

// GetConnectionRepository mengembalikan repository koneksi
func (c *Container) GetConnectionRepository() repository.ConnectionRepository {
	return c.connectionRepository
}

// GetCommandRepository mengembalikan repository command
func (c *Container) GetCommandRepository() repository.CommandRepository {
	return c.commandRepository
}

// GetStatsRepository mengembalikan repository statistik
func (c *Container) GetStatsRepository() repository.StatsRepository {
	return c.statsRepository
}

// GetFinanceService mengembalikan finance service
func (c *Container) GetFinanceService() service.FinanceService {
	return c.financeService
}

// GetGoogleAPIRepository mengembalikan repository Google API
func (c *Container) GetGoogleAPIRepository() repository.GoogleAPIRepository {
	return c.googleAPIRepository
}

// GetDriveRepository mengembalikan repository drive
func (c *Container) GetDriveRepository() repository.DriveRepository {
	return c.driveRepository
}

// GetConnectWhatsAppUseCase mengembalikan use case koneksi WhatsApp
func (c *Container) GetConnectWhatsAppUseCase() *connectionUseCase.ConnectWhatsAppUseCase {
	return c.connectWhatsAppUseCase
}

// GetDashboardController mengembalikan controller dashboard
func (c *Container) GetDashboardController() *web.DashboardController {
	return c.dashboardController
}

// GetQRController mengembalikan controller QR
func (c *Container) GetQRController() *web.QRController {
	return c.qrController
}

// GetAuthController mengembalikan controller autentikasi
func (c *Container) GetAuthController() *web.AuthController {
	return c.authController
}

// GetMessageController mengembalikan controller pesan
func (c *Container) GetMessageController() *whatsappController.MessageController {
	return c.messageController
}

// GetConfigController mengembalikan controller konfigurasi
func (c *Container) GetConfigController() *web.ConfigController {
	return c.configController
}

// GetDataMasterController mengembalikan controller data master
func (c *Container) GetDataMasterController() *web.DataMasterController {
	return c.dataMasterController
}

// GetContactService mengembalikan service kontak
func (c *Container) GetContactService() service.ContactService {
	return c.contactService
}

// GetContactController mengembalikan controller kontak
func (c *Container) GetContactController() interface{} {
	return c.contactController
}

// GetConfig mengembalikan konfigurasi
func (c *Container) GetConfig() *config.Config {
	return c.config
}

// GetPort mengembalikan port dari konfigurasi
func (c *Container) GetPort() int {
	return c.config.GetWebPort()
}

// whatsmeowLoggerAdapter adalah adapter untuk logger whatsmeow
type whatsmeowLoggerAdapter struct {
	log *logger.Logger
}

func (l *whatsmeowLoggerAdapter) Debugf(format string, args ...interface{}) {
	l.log.Debug(format, args...)
}

func (l *whatsmeowLoggerAdapter) Infof(format string, args ...interface{}) {
	l.log.Info(format, args...)
}

func (l *whatsmeowLoggerAdapter) Warnf(format string, args ...interface{}) {
	l.log.Warn(format, args...)
}

func (l *whatsmeowLoggerAdapter) Errorf(format string, args ...interface{}) {
	l.log.Error(format, args...)
}

func (l *whatsmeowLoggerAdapter) Sub(module string) waLog.Logger {
	return &whatsmeowLoggerAdapter{
		log: logger.New("WhatsmeowDB."+module, logger.INFO, true),
	}
}
