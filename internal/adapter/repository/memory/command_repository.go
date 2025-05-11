package memory

import (
	"sync"

	"github.com/gwenziro/botopia/internal/domain/command"
)

// CommandRepository implementasi in-memory untuk penyimpanan command
type CommandRepository struct {
	commands map[string]command.Command
	mutex    sync.RWMutex
}

// NewCommandRepository membuat instance repository command baru
func NewCommandRepository() *CommandRepository {
	return &CommandRepository{
		commands: make(map[string]command.Command),
	}
}

// FindByName mencari command berdasarkan nama
func (r *CommandRepository) FindByName(name string) (command.Command, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	cmd, found := r.commands[name]
	return cmd, found
}

// GetAll mengembalikan semua command
func (r *CommandRepository) GetAll() map[string]command.Command {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Buat salinan untuk menghindari race condition
	result := make(map[string]command.Command, len(r.commands))
	for name, cmd := range r.commands {
		result[name] = cmd
	}

	return result
}

// Register mendaftarkan command baru
func (r *CommandRepository) Register(cmd command.Command) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.commands[cmd.GetName()] = cmd
	return nil
}

// GetByCategory mengembalikan command berdasarkan kategori
func (r *CommandRepository) GetByCategory(category string) []command.Command {
	// Implementasi sederhana, bisa dikembangkan nanti untuk mendukung kategori
	return nil
}
