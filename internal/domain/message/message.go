package message

import (
	"fmt"
	"strings"
	"time"

	"github.com/gwenziro/botopia/internal/domain/user"
)

// MediaType mendefinisikan tipe media yang didukung
type MediaType string

const (
	// MediaTypeImage adalah tipe media gambar
	MediaTypeImage MediaType = "image"

	// MediaTypeVideo adalah tipe media video
	MediaTypeVideo MediaType = "video"

	// MediaTypeDocument adalah tipe media dokumen
	MediaTypeDocument MediaType = "document"
)

// Chat merepresentasikan chat WhatsApp
type Chat struct {
	ID      string
	IsGroup bool
	Name    string
}

// Message merepresentasikan pesan WhatsApp
type Message struct {
	ID         string
	Text       string
	Sender     *user.User
	Chat       *Chat
	Timestamp  time.Time
	IsFromMe   bool
	IsGroup    bool
	MediaType  MediaType
	Caption    string
	MediaURL   string
	ReplyingTo *Message
}

// HasMedia memeriksa apakah pesan mengandung media
func (m *Message) HasMedia() bool {
	return m.MediaType != ""
}

// ExtractCommand mengekstrak nama command dan argumen dari teks pesan
func (m *Message) ExtractCommand(prefix string) (string, []string, bool) {
	// Cek pesan teks
	var textToCheck string

	if m.Text != "" {
		textToCheck = m.Text
	} else if m.Caption != "" {
		// Jika teks kosong tapi caption ada (pesan media), gunakan caption
		textToCheck = m.Caption
	}

	if textToCheck == "" {
		return "", nil, false
	}

	text := strings.TrimSpace(textToCheck)

	// Periksa apakah pesan dimulai dengan prefix command
	if !strings.HasPrefix(text, prefix) {
		return "", nil, false
	}

	// Hapus prefix
	text = text[len(prefix):]

	// Split text menjadi command dan args
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return "", nil, false
	}

	cmdName := strings.ToLower(parts[0])
	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	} else {
		args = []string{}
	}

	return cmdName, args, true
}

// DownloadMedia mengunduh media dari pesan ini dan mengembalikan path file lokal
func (m *Message) DownloadMedia() (string, error) {
	if !m.HasMedia() {
		return "", fmt.Errorf("pesan tidak memiliki media")
	}

	// Gunakan singleton untuk mendapatkan layanan
	if downloadMediaService == nil {
		return "", fmt.Errorf("layanan unduh media belum diinisialisasi")
	}

	return downloadMediaService.DownloadMedia(m)
}

// DownloadMediaService antarmuka untuk mengunduh media
type DownloadMediaService interface {
	DownloadMedia(msg *Message) (string, error)
}

var downloadMediaService DownloadMediaService

// SetDownloadMediaService menetapkan layanan untuk unduh media
func SetDownloadMediaService(service DownloadMediaService) {
	downloadMediaService = service
}
