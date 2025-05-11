package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config menyimpan konfigurasi aplikasi
type Config struct {
	// Database
	DBPath   string
	DBBackup bool

	// Bot
	CommandPrefix string

	// Logging
	LogLevel  string
	UseColors bool

	// Web
	WebPort         int
	WebHost         string
	WebAuthEnabled  bool
	WebAuthUsername string
	WebAuthPassword string

	// Aplikasi
	DevMode       bool
	CleanStart    bool
	QRAutoRefresh bool

	// Direktori
	WebViewDir   string
	WebStaticDir string
}

// NewConfig membuat instance Config baru dengan nilai default
func NewConfig() *Config {
	return &Config{
		DBPath:          "botopia.db",
		DBBackup:        true,
		CommandPrefix:   "!",
		LogLevel:        "INFO",
		UseColors:       true,
		WebPort:         8080,
		WebHost:         "0.0.0.0",
		WebAuthEnabled:  false,
		WebAuthUsername: "admin",
		WebAuthPassword: "admin",
		DevMode:         false,
		CleanStart:      false,
		QRAutoRefresh:   true,
		WebViewDir:      "./internal/infrastructure/web/view",
		WebStaticDir:    "./internal/infrastructure/web/static",
	}
}

// LoadFromEnv memuat konfigurasi dari environment variables
func (c *Config) LoadFromEnv() {
	// Database
	if v := os.Getenv("BOTOPIA_DB_PATH"); v != "" {
		c.DBPath = v
	}

	if v := os.Getenv("BOTOPIA_DB_BACKUP"); v != "" {
		c.DBBackup = strings.ToLower(v) == "true"
	}

	// Bot
	if v := os.Getenv("BOTOPIA_COMMAND_PREFIX"); v != "" {
		c.CommandPrefix = v
	}

	// Logging
	if v := os.Getenv("BOTOPIA_LOG_LEVEL"); v != "" {
		c.LogLevel = v
	}

	if v := os.Getenv("BOTOPIA_USE_COLORS"); v != "" {
		c.UseColors = strings.ToLower(v) == "true"
	}

	// Web
	if v := os.Getenv("BOTOPIA_WEB_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			c.WebPort = port
		}
	}

	if v := os.Getenv("BOTOPIA_WEB_HOST"); v != "" {
		c.WebHost = v
	}

	if v := os.Getenv("BOTOPIA_AUTH_ENABLED"); v != "" {
		c.WebAuthEnabled = strings.ToLower(v) == "true"
	}

	if v := os.Getenv("BOTOPIA_AUTH_USERNAME"); v != "" {
		c.WebAuthUsername = v
	}

	if v := os.Getenv("BOTOPIA_AUTH_PASSWORD"); v != "" {
		c.WebAuthPassword = v
	}

	// Aplikasi
	if v := os.Getenv("BOTOPIA_DEV_MODE"); v != "" {
		c.DevMode = strings.ToLower(v) == "true"
	}

	if v := os.Getenv("BOTOPIA_CLEAN_START"); v != "" {
		c.CleanStart = strings.ToLower(v) == "true"
	}

	if v := os.Getenv("BOTOPIA_QR_AUTO_REFRESH"); v != "" {
		c.QRAutoRefresh = strings.ToLower(v) == "true"
	}

	// Cek direktori view dan static baru
	if _, err := os.Stat("./internal/infrastructure/web/view"); err == nil {
		c.WebViewDir = "./internal/infrastructure/web/view"
	}

	if _, err := os.Stat("./internal/infrastructure/web/static"); err == nil {
		c.WebStaticDir = "./internal/infrastructure/web/static"
	}

	// Cek custom direktori dari env
	if v := os.Getenv("BOTOPIA_WEB_VIEW_DIR"); v != "" {
		if _, err := os.Stat(v); err == nil {
			c.WebViewDir = v
		}
	}

	if v := os.Getenv("BOTOPIA_WEB_STATIC_DIR"); v != "" {
		if _, err := os.Stat(v); err == nil {
			c.WebStaticDir = v
		}
	}
}

// GetWebPort mengembalikan port web server
func (c *Config) GetWebPort() int {
	return c.WebPort
}

// IsDevMode memeriksa apakah dalam mode pengembangan
func (c *Config) IsDevMode() bool {
	return c.DevMode
}

// EnsureDirectories memastikan semua direktori yang diperlukan tersedia
func (c *Config) EnsureDirectories() error {
	dirs := []string{
		filepath.Dir(c.DBPath),
		c.WebViewDir,
		c.WebStaticDir,
	}

	for _, dir := range dirs {
		if dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}
	}

	return nil
}
