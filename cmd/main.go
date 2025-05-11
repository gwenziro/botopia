package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gwenziro/botopia/internal/app/di"
	"github.com/gwenziro/botopia/internal/infrastructure/config"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
	"github.com/gwenziro/botopia/internal/infrastructure/server"
)

func main() {
	// Inisialisasi konfigurasi
	cfg := config.NewConfig()
	cfg.LoadFromEnv()

	// Setup logger
	log := logger.New("Main", logger.LevelFromString(cfg.LogLevel), cfg.UseColors)
	log.Info("Starting Botopia with configuration: DevMode=%v, WebPort=%d", cfg.IsDevMode(), cfg.GetWebPort())

	// Inisialisasi dependency injection container
	log.Info("Initializing DI container...")
	container := di.NewContainer(cfg)

	// Start web server
	log.Info("Starting web server...")
	webServer := server.NewServer(container)
	go func() {
		if err := webServer.Start(); err != nil {
			log.Fatal("Error starting web server: %v", err)
		}
	}()

	// Connect to WhatsApp
	log.Info("Connecting to WhatsApp...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	status, err := container.GetConnectWhatsAppUseCase().Execute(ctx)
	if err != nil {
		log.Warn("Failed to connect to WhatsApp: %v", err)
		log.Info("Please scan QR code to connect")
	} else {
		log.Info("WhatsApp connection status: %v", status.IsConnected)
	}

	// Setup signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Botopia is running. Press Ctrl+C to exit.")

	// Wait for interrupt signal
	<-quit
	log.Info("Shutting down...")

	// Create context with timeout for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Shutdown web server
	if err := webServer.Shutdown(shutdownCtx); err != nil {
		log.Error("Error shutting down web server: %v", err)
	}

	// Disconnect WhatsApp
	log.Info("Disconnecting from WhatsApp...")
	container.GetConnectionRepository().Disconnect()

	log.Info("Server gracefully stopped")
}
