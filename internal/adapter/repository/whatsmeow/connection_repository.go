package whatsmeow

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

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

// Maps untuk menyimpan referensi pesan untuk akses media
var (
	messageCache      = map[string]*events.Message{}
	messageCacheMutex sync.RWMutex
)

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

// DownloadMedia mengunduh media dari pesan
func (r *ConnectionRepository) DownloadMedia(ctx context.Context, msg *message.Message) (string, error) {
	r.log.Info("Mengunduh media dari pesan: %s", msg.ID)

	// Buat temporary directory jika belum ada
	tempDir := "./tmp/media"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("gagal membuat direktori temporary: %v", err)
	}

	// Tentukan nama file dan path
	timestamp := time.Now().Unix()
	fileExt := "jpg" // Default untuk gambar
	switch msg.MediaType {
	case message.MediaTypeVideo:
		fileExt = "mp4"
	case message.MediaTypeDocument:
		fileExt = "pdf" // Default for documents
	}

	fileName := fmt.Sprintf("%s/proof_%d.%s", tempDir, timestamp, fileExt)

	// Kita perlu mengekstrak info media dari message.Message asli (WhatsApp) melalui EventInfo
	// Karena whatsmeow tidak menyediakan GetMessage, kita perlu menyimpan informasi media
	// saat menerima pesan asli

	// Coba ambil dari map cache messages jika implementasi EventMessage original tersedia
	var downloadable whatsmeow.DownloadableMessage
	_, err := waTypes.ParseJID(msg.Chat.ID)
	if err != nil {
		r.log.Error("Gagal parse chat JID: %v", err)
		return "", fmt.Errorf("gagal parse chat JID: %v", err)
	}

	// Cari sumber media berdasarkan jenis pesan
	switch msg.MediaType {
	case message.MediaTypeImage:
		// Coba dapatkan image message dari store cache
		evt, err := r.getMessageFromStore(msg)
		if err != nil {
			return "", err
		}
		if evt.Message.ImageMessage != nil {
			downloadable = evt.Message.ImageMessage
		}
	case message.MediaTypeVideo:
		// Coba dapatkan video message dari store cache
		evt, err := r.getMessageFromStore(msg)
		if err != nil {
			return "", err
		}
		if evt.Message.VideoMessage != nil {
			downloadable = evt.Message.VideoMessage
		}
	case message.MediaTypeDocument:
		// Coba dapatkan document message dari store cache
		evt, err := r.getMessageFromStore(msg)
		if err != nil {
			return "", err
		}
		if evt.Message.DocumentMessage != nil {
			downloadable = evt.Message.DocumentMessage
		}
	default:
		return "", fmt.Errorf("tipe media tidak didukung: %s", string(msg.MediaType))
	}

	if downloadable == nil {
		return "", fmt.Errorf("tidak dapat menemukan media dalam pesan")
	}

	// Unduh media
	mediaData, err := r.client.Download(downloadable)
	if err != nil {
		r.log.Error("Gagal mengunduh media: %v", err)
		return "", fmt.Errorf("gagal mengunduh media: %v", err)
	}

	// Tulis file ke disk
	if err := os.WriteFile(fileName, mediaData, 0644); err != nil {
		r.log.Error("Gagal menyimpan media: %v", err)
		return "", fmt.Errorf("gagal menyimpan media: %v", err)
	}

	r.log.Info("Media berhasil diunduh ke: %s", fileName)
	return fileName, nil
}

// getMessageFromStore mencoba mendapatkan pesan asli dari store berdasarkan ID
func (r *ConnectionRepository) getMessageFromStore(msg *message.Message) (*events.Message, error) {
	// Log untuk debugging
	r.log.Debug("Mencari pesan dengan ID: %s di cache", msg.ID)

	// Cek apakah pesan ada di cache
	messageCacheMutex.RLock()
	evt, exists := messageCache[msg.ID]
	messageCacheMutex.RUnlock()

	if exists {
		r.log.Debug("Pesan ditemukan di cache")
		return evt, nil
	}

	r.log.Warn("Pesan dengan ID %s tidak ditemukan dalam cache", msg.ID)

	// Coba cari berdasarkan ID dan timestamp untuk antisipasi
	var foundEvt *events.Message
	messageCacheMutex.RLock()
	for id, cachedEvt := range messageCache {
		if cachedEvt.Info.Timestamp.Unix() == msg.Timestamp.Unix() {
			r.log.Info("Menemukan pesan berdasarkan timestamp, ID: %s", id)
			foundEvt = cachedEvt
			break
		}
	}
	messageCacheMutex.RUnlock()

	if foundEvt != nil {
		return foundEvt, nil
	}

	// Jika tidak ditemukan sama sekali, kembalikan error
	return nil, fmt.Errorf("pesan dengan ID %s tidak ditemukan dalam cache", msg.ID)
}

// handleEvent menangani event dari WhatsApp
func (r *ConnectionRepository) handleEvent(evt interface{}) {
	if evt == nil {
		return
	}

	switch v := evt.(type) {
	case *events.Message:
		// Log detail pesan untuk debugging
		if v.Message.ImageMessage != nil {
			r.log.Debug("Menerima pesan gambar")
			if v.Message.ImageMessage.Caption != nil {
				r.log.Debug("Caption gambar: %s", *v.Message.ImageMessage.Caption)
			}
		} else if v.Message.VideoMessage != nil {
			r.log.Debug("Menerima pesan video")
		} else if v.Message.DocumentMessage != nil {
			r.log.Debug("Menerima pesan dokumen")
		}

		// Simpan pesan media di cache sementara untuk akses nanti
		if v.Message.ImageMessage != nil || v.Message.VideoMessage != nil || v.Message.DocumentMessage != nil {
			r.log.Debug("Menyimpan pesan media dengan ID: %s di cache", v.Info.ID)
			messageCacheMutex.Lock()
			messageCache[v.Info.ID] = v
			messageCacheMutex.Unlock()

			// Set goroutine untuk membersihkan cache setelah beberapa waktu
			go func(id string) {
				time.Sleep(30 * time.Minute) // Cache pesan untuk 30 menit
				messageCacheMutex.Lock()
				delete(messageCache, id)
				messageCacheMutex.Unlock()
			}(v.Info.ID)
		}

		// Konversi ke domain message
		msg := r.convertToMessage(v)

		// Log untuk debugging
		if msg.HasMedia() {
			r.log.Info("Pesan media berhasil dikonversi, ID: %s, MediaType: %s, Caption: %s",
				msg.ID, string(msg.MediaType), msg.Caption)
		}

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

	// Tentukan jenis media dan caption
	var mediaType message.MediaType
	var caption string

	if evt.Message.ImageMessage != nil {
		mediaType = message.MediaTypeImage
		if evt.Message.ImageMessage.Caption != nil {
			caption = *evt.Message.ImageMessage.Caption
			// Log untuk debugging
			r.log.Debug("Caption dari image: %s", caption)
		}
	} else if evt.Message.VideoMessage != nil {
		mediaType = message.MediaTypeVideo
		if evt.Message.VideoMessage.Caption != nil {
			caption = *evt.Message.VideoMessage.Caption
		}
	} else if evt.Message.DocumentMessage != nil {
		mediaType = message.MediaTypeDocument
		if evt.Message.DocumentMessage.Title != nil {
			caption = *evt.Message.DocumentMessage.Title
		}
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

	msg := &message.Message{
		ID:        evt.Info.ID,
		Text:      text,
		Sender:    u,
		Chat:      chat,
		Timestamp: evt.Info.Timestamp,
		IsFromMe:  evt.Info.IsFromMe,
		IsGroup:   evt.Info.IsGroup,
		MediaType: mediaType,
		Caption:   caption,
	}

	// Log detail pesan media untuk debugging
	if mediaType != "" {
		r.log.Debug("Konversi message dengan media type: %s, caption: %s", string(mediaType), caption)
	}

	return msg
}
