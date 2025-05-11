package ping

import (
	"testing"
	"time"

	"github.com/gwenziro/botopia/internal/domain/message"
	"github.com/gwenziro/botopia/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

func TestPingCommand_GetName(t *testing.T) {
	// Arrange
	cmd := NewCommand()

	// Act
	name := cmd.GetName()

	// Assert
	assert.Equal(t, "ping", name)
}

func TestPingCommand_Execute(t *testing.T) {
	// Arrange
	cmd := NewCommand()
	msg := &message.Message{
		ID:   "test-id",
		Text: "!ping",
		Sender: &user.User{
			ID:    "123",
			Phone: "+123456789",
		},
		Timestamp: time.Now().Add(-100 * time.Millisecond),
	}

	// Act
	response, err := cmd.Execute([]string{}, msg)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, response, "Pong!")
	assert.Contains(t, response, "Bot sedang aktif")
}
