package ping

import (
	"fmt"
	"time"

	"github.com/gwenziro/botopia/internal/domain/message"
)

// Command implementasi command ping
type Command struct{}

// NewCommand membuat instance ping command baru
func NewCommand() *Command {
	return &Command{}
}

// GetName mengembalikan nama command
func (c *Command) GetName() string {
	return "ping"
}

// GetDescription mengembalikan deskripsi command
func (c *Command) GetDescription() string {
	return "Mengecek apakah bot sedang aktif dan menampilkan waktu respons"
}

// Execute menjalankan command
func (c *Command) Execute(args []string, msg *message.Message) (string, error) {
	responseTime := time.Since(msg.Timestamp)
	return fmt.Sprintf("Pong! Bot sedang aktif. Waktu respons: %s", responseTime.String()), nil
}
