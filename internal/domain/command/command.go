package command

import (
	"github.com/gwenziro/botopia/internal/domain/message"
)

// Command mendefinisikan kontrak untuk command bot
type Command interface {
	// GetName mengembalikan nama command
	GetName() string

	// GetDescription mengembalikan deskripsi command
	GetDescription() string

	// Execute menjalankan command dan mengembalikan response
	Execute(args []string, msg *message.Message) (string, error)

	// GetCategory mengembalikan kategori command (opsional)
	GetCategory() string

	// GetUsage mengembalikan contoh penggunaan command (opsional)
	GetUsage() string
}

// CommandInfo berisi informasi tentang command
type CommandInfo struct {
	Name        string
	Description string
	Usage       string
	Category    string
}
