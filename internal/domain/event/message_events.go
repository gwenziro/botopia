package event

import (
	"github.com/gwenziro/botopia/internal/domain/message"
)

// MessageReceivedEvent adalah event ketika pesan diterima
type MessageReceivedEvent struct {
	BaseEvent
	Message *message.Message
}

// NewMessageReceivedEvent membuat instance MessageReceivedEvent baru
func NewMessageReceivedEvent(msg *message.Message) MessageReceivedEvent {
	return MessageReceivedEvent{
		BaseEvent: NewBaseEvent("message.received"),
		Message:   msg,
	}
}

// CommandExecutedEvent adalah event ketika command dijalankan
type CommandExecutedEvent struct {
	BaseEvent
	CommandName string
	Message     *message.Message
	Response    string
}

// NewCommandExecutedEvent membuat instance CommandExecutedEvent baru
func NewCommandExecutedEvent(cmdName string, msg *message.Message, response string) CommandExecutedEvent {
	return CommandExecutedEvent{
		BaseEvent:   NewBaseEvent("command.executed"),
		CommandName: cmdName,
		Message:     msg,
		Response:    response,
	}
}
