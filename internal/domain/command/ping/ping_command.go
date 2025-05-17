package ping

import (
	"fmt"
	"time"

	"github.com/gwenziro/botopia/internal/domain/command/common"
	"github.com/gwenziro/botopia/internal/domain/message"
)

// Command implementasi command ping
type Command struct {
	common.BaseCommand
}

// NewCommand membuat instance ping command baru
func NewCommand() *Command {
	cmd := &Command{}
	cmd.Name = "ping"
	cmd.Description = "Mengecek apakah bot sedang aktif dan menampilkan waktu respons"
	cmd.Category = "System"
	cmd.Usage = "!ping"
	return cmd
}

// Execute menjalankan command
func (c *Command) Execute(args []string, msg *message.Message) (string, error) {
	responseTime := time.Since(msg.Timestamp)
	return fmt.Sprintf("Pong! Bot sedang aktif. Waktu respons: %s", responseTime.String()), nil
}
