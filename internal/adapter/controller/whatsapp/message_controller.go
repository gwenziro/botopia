package whatsapp

import (
	"context"
	"time"

	"github.com/gwenziro/botopia/internal/domain/event"
	"github.com/gwenziro/botopia/internal/domain/message"
	"github.com/gwenziro/botopia/internal/domain/repository"
	"github.com/gwenziro/botopia/internal/domain/service"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
	"github.com/gwenziro/botopia/internal/usecase/command/execute"
)

// MessageController controller untuk pesan WhatsApp
type MessageController struct {
	executeCommandUseCase *execute.ExecuteCommandUseCase
	connectionRepo        repository.ConnectionRepository
	statsRepo             repository.StatsRepository
	contactService        service.ContactService
	log                   *logger.Logger
	eventDispatcher       *event.EventDispatcher
	UseWhitelist          bool // Diubah menjadi exported (huruf kapital)
}

// NewMessageController membuat instance message controller
func NewMessageController(
	executeUC *execute.ExecuteCommandUseCase,
	connectionRepo repository.ConnectionRepository,
	statsRepo repository.StatsRepository,
	contactService service.ContactService,
) *MessageController {
	return &MessageController{
		executeCommandUseCase: executeUC,
		connectionRepo:        connectionRepo,
		statsRepo:             statsRepo,
		contactService:        contactService,
		log:                   logger.New("MessageController", logger.INFO, true),
		eventDispatcher:       event.NewEventDispatcher(),
		UseWhitelist:          false, // Default: nonaktif
	}
}

// SetUseWhitelist mengatur apakah menggunakan whitelist
func (c *MessageController) SetUseWhitelist(value bool) {
	c.log.Info("Pengaturan whitelist diubah: %v", value)
	c.UseWhitelist = value
}

// Setup menyiapkan controller
func (c *MessageController) Setup() {
	// Daftarkan handler pesan
	c.connectionRepo.RegisterMessageHandler(c.HandleMessage)
	c.log.Info("Message handler registered")
}

// HandleMessage menangani pesan masuk dengan validasi whitelist
func (c *MessageController) HandleMessage(msg *message.Message) {
	// Increment pesan diterima
	c.statsRepo.IncrementMessageCount()

	// Dispatch message received event
	c.eventDispatcher.Dispatch(event.NewMessageReceivedEvent(msg))

	// Abaikan pesan yang dikirim oleh kita
	if msg.IsFromMe {
		return
	}

	// Validasi whitelist jika diaktifkan
	if c.UseWhitelist {
		// Cek apakah pengirim ada dalam whitelist
		senderPhone := msg.Sender.Phone
		isWhitelisted, err := c.contactService.IsWhitelisted(context.Background(), senderPhone)

		if err != nil {
			c.log.Error("Gagal memeriksa whitelist untuk %s: %v", senderPhone, err)
			// Lanjutkan proses pesannya, jangan blokir karena error teknis
		} else if !isWhitelisted {
			c.log.Info("Pesan dari nomor yang tidak dalam whitelist ditolak: %s", senderPhone)
			// Opsional: kirim respons bahwa nomor tidak diizinkan
			// c.sendReply(msg, "Maaf, nomor Anda tidak terdaftar untuk menggunakan layanan ini.")
			return
		}
	}

	// Gunakan caption sebagai text jika ada media dan caption
	if msg.HasMedia() && msg.Caption != "" {
		// Tambahkan log untuk debug
		c.log.Info("Pesan dengan media diterima, caption: %s", msg.Caption)
		msg.Text = msg.Caption
	} else {
		// Log jika pesan media tanpa caption
		if msg.HasMedia() {
			c.log.Info("Pesan dengan media diterima tanpa caption")
		}
	}

	// Log pesan masuk
	c.log.Debug("Message received: %s", msg.Text)

	// Buat context dengan timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Jalankan command jika ada
	response, err := c.executeCommandUseCase.Execute(ctx, msg)
	if err != nil {
		if err == execute.ErrCommandNotFound {
			// Command tidak ditemukan
			if msg.HasMedia() {
				c.log.Info("Media diterima tanpa command yang valid")
				c.sendReply(msg, "Untuk mengunggah bukti transaksi, gunakan format caption: !unggah <kode_transaksi>")
			}
			return
		}
		c.log.Error("Failed to execute command: %v", err)
		return
	}

	// Dispatch command executed event
	if response != "" {
		cmdName, _, _ := msg.ExtractCommand("!")
		c.eventDispatcher.Dispatch(event.NewCommandExecutedEvent(cmdName, msg, response))
		c.sendReply(msg, response)
	}
}

// sendReply mengirim balasan ke pesan
func (c *MessageController) sendReply(msg *message.Message, text string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.connectionRepo.SendMessage(ctx, msg.Chat.ID, text)
	if err != nil {
		c.log.Error("Failed to send reply: %v", err)
	}
}
