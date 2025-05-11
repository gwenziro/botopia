package repository

import "github.com/gwenziro/botopia/internal/domain/command"

// CommandRepository mendefinisikan kontrak untuk akses data command
type CommandRepository interface {
	// FindByName mencari command berdasarkan nama
	FindByName(name string) (command.Command, bool)

	// GetAll mengembalikan semua command
	GetAll() map[string]command.Command

	// Register mendaftarkan command baru
	Register(cmd command.Command) error

	// GetByCategory mengembalikan command berdasarkan kategori
	GetByCategory(category string) []command.Command
}
