package whatsmeow

import (
	"context"
	"sync"

	"github.com/gwenziro/botopia/internal/domain/message"
	"github.com/gwenziro/botopia/internal/domain/user"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

// ConnectionRepository implementasi repository koneksi WhatsApp
type ConnectionRepository struct {
	client          *whatsmeow.Client
	log             *logger.Logger
	isConnected     bool
	connMutex       sync.RWMutex
	qrChan          chan string
	messageHandlers []func(*message.Message)
	handlersMutex   sync.RWMutex
}

// NewConnectionRepository membuat repository koneksi baru
func NewConnectionRepository(client *whatsmeow.Client, log *logger.Logger) *ConnectionRepository {
	repo := &ConnectionRepository{
		client:          client,
		log:             log,
		qrChan:          make(chan string, 5),
		messageHandlers: make([]func(*message.Message), 0),
	}

	// Daftarkan event handler
	client.AddEventHandler(repo.handleEvent)

	return repo
}

// Connect menghubungkan ke WhatsApp
func (r *ConnectionRepository) Connect(ctx context.Context) error {
	r.connMutex.Lock()
	defer r.connMutex.Unlock()

	if r.client.Store.ID == nil {
		// Baru, butuh QR code
		qrChan, _ := r.client.GetQRChannel(ctx)
		err := r.client.Connect()
		if err != nil {
			return err
		}

		// Forward QR code events
		go func() {
			for evt := range qrChan {
				if evt.Event == "code" {
					r.qrChan <- evt.Code
				}
			}
		}()
	} else {
		// Gunakan session yang ada
		err := r.client.Connect()
		if err != nil {
			return err
		}
	}

	return nil
}

// Disconnect memutuskan koneksi WhatsApp
func (r *ConnectionRepository) Disconnect() error {
	r.connMutex.Lock()
	defer r.connMutex.Unlock()

	r.client.Disconnect()
	r.isConnected = false

	return nil
}

// IsConnected memeriksa status koneksi
func (r *ConnectionRepository) IsConnected() bool {
	r.connMutex.RLock()
	defer r.connMutex.RUnlock()

	return r.isConnected
}

// GetCurrentUser mendapatkan informasi user yang terhubung
func (r *ConnectionRepository) GetCurrentUser() (*user.User, error) {
	if !r.IsConnected() || r.client == nil || r.client.Store == nil || r.client.Store.ID == nil {
		return nil, nil
	}

	// Ekstrak informasi user dari client
	jid := r.client.Store.ID
	if jid == nil {
		return nil, nil
	}

	return &user.User{
		ID:    jid.User,
		Phone: "+" + jid.User,
		Name:  "WhatsApp User", // Default name
	}, nil
}

// SendMessage mengirim pesan
func (r *ConnectionRepository) SendMessage(ctx context.Context, chatID string, text string) error {
	jid, err := waTypes.ParseJID(chatID)
	if err != nil {
		return err
	}

	_, err = r.client.SendMessage(ctx, jid, &waProto.Message{
		Conversation: proto.String(text),
	})

	return err
}

// RegisterMessageHandler mendaftarkan handler untuk pesan masuk
func (r *ConnectionRepository) RegisterMessageHandler(handler func(*message.Message)) {
	r.handlersMutex.Lock()
	defer r.handlersMutex.Unlock()

	r.messageHandlers = append(r.messageHandlers, handler)
}

// GetQRChannel mendapatkan channel untuk QR code
func (r *ConnectionRepository) GetQRChannel() <-chan string {
	return r.qrChan
}

// handleEvent menangani event dari WhatsApp
func (r *ConnectionRepository) handleEvent(evt interface{}) {
	if evt == nil {
		return
	}

	switch v := evt.(type) {
	case *events.Message:
		msg := r.convertToMessage(v)

		r.handlersMutex.RLock()
		handlers := r.messageHandlers
		r.handlersMutex.RUnlock()

		for _, handler := range handlers {
			go handler(msg)
		}

	case *events.Connected:
		r.connMutex.Lock()
		r.isConnected = true
		r.connMutex.Unlock()
		r.log.Info("WhatsApp terhubung!")

	case *events.Disconnected:
		r.connMutex.Lock()
		r.isConnected = false
		r.connMutex.Unlock()
		r.log.Info("WhatsApp terputus")
	}
}

// convertToMessage mengubah event message ke domain message
func (r *ConnectionRepository) convertToMessage(evt *events.Message) *message.Message {
	// Ekstrak text message
	text := ""
	if evt.Message.GetConversation() != "" {
		text = evt.Message.GetConversation()
	} else if evt.Message.ExtendedTextMessage != nil && evt.Message.ExtendedTextMessage.Text != nil {
		text = *evt.Message.ExtendedTextMessage.Text
	}

	// Buat user object
	u := &user.User{
		ID:    evt.Info.Sender.User,
		Phone: "+" + evt.Info.Sender.User,
	}

	// Buat chat object
	chat := &message.Chat{
		ID:      evt.Info.Chat.String(),
		IsGroup: evt.Info.IsGroup,
	}

	return &message.Message{
		ID:        evt.Info.ID,
		Text:      text,
		Sender:    u,
		Chat:      chat,
		Timestamp: evt.Info.Timestamp,
		IsFromMe:  evt.Info.IsFromMe,
		IsGroup:   evt.Info.IsGroup,
	}
}
