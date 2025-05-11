package message

import (
	"strings"
	"time"

	"github.com/gwenziro/botopia/internal/domain/user"
)

// Message merepresentasikan pesan WhatsApp
type Message struct {
	ID        string
	Text      string
	Sender    *user.User
	Chat      *Chat
	Timestamp time.Time
	IsFromMe  bool
	IsGroup   bool
}

// Chat merepresentasikan chat WhatsApp
type Chat struct {
	ID      string
	Name    string
	IsGroup bool
}

// ExtractCommand mengekstrak command dari pesan jika ada
func (m *Message) ExtractCommand(prefix string) (cmdName string, args []string, isCommand bool) {
	if !strings.HasPrefix(m.Text, prefix) {
		return "", nil, false
	}

	// Parse command dan arguments
	commandText := strings.TrimPrefix(m.Text, prefix)
	parts := strings.Fields(commandText)

	if len(parts) == 0 {
		return "", nil, false
	}

	cmdName = parts[0]
	args = []string{}

	if len(parts) > 1 {
		args = parts[1:]
	}

	return cmdName, args, true
}
