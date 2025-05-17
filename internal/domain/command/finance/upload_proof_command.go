package finance

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gwenziro/botopia/internal/domain/command/common"
	"github.com/gwenziro/botopia/internal/domain/dto"
	"github.com/gwenziro/botopia/internal/domain/finance"
	"github.com/gwenziro/botopia/internal/domain/message"
	"github.com/gwenziro/botopia/internal/domain/service"
)

// UploadProofCommand implementasi command untuk mengunggah bukti transaksi
type UploadProofCommand struct {
	common.BaseCommand
	financeService service.FinanceService
}

// NewUploadProofCommand membuat instance command baru
func NewUploadProofCommand(financeService service.FinanceService) *UploadProofCommand {
	cmd := &UploadProofCommand{
		financeService: financeService,
	}
	cmd.Name = "unggah"
	cmd.Description = "Mengunggah bukti transaksi untuk catatan yang sudah ada. Kirim !unggah untuk mendapatkan form input data."
	cmd.Category = "Keuangan"
	cmd.Usage = "!unggah <kode_transaksi>"
	return cmd
}

// Execute menjalankan command
func (c *UploadProofCommand) Execute(args []string, msg *message.Message) (string, error) {
	// Jika tidak ada argumen, kirimkan form template
	if len(args) == 0 {
		return c.getFormTemplate(), nil
	}

	// Cek apakah pesan adalah form yang diisi
	if isFilledForm, form := c.parseFormInput(msg.Text); isFilledForm {
		// Validasi media
		if !msg.HasMedia() {
			return "❌ Mohon lampirkan foto bukti transaksi untuk diunggah.", nil
		}

		// Ambil kode transaksi dari form
		transactionCode := form["Kode unik"]
		return c.processUpload(msg, transactionCode)
	}

	// Jika format tidak sesuai, berikan panduan penggunaan form
	return "❌ Format tidak sesuai. Silakan gunakan format formulir yang tersedia dengan mengirim !unggah tanpa parameter tambahan.", nil
}

// getFormTemplate mengembalikan template form unggah bukti
func (c *UploadProofCommand) getFormTemplate() string {
	return `!unggah
────────────────────────
ℹ FORMAT UNGGAH BUKTI TRANSAKSI ℹ
────────────────────────
Kode unik: 

────────────────────────
Tolong kirim bukti transaksi dengan format di atas, ya! 🙏`
}

// parseFormInput mengecek dan mem-parsing input formulir
func (c *UploadProofCommand) parseFormInput(text string) (bool, map[string]string) {
	// Cek apakah dimulai dengan !unggah
	if !strings.HasPrefix(text, "!unggah") {
		return false, nil
	}

	// Regex untuk mengekstrak nilai dari field
	form := make(map[string]string)

	// Ekstrak kode unik
	re := regexp.MustCompile(`Kode unik:\s*(.*?)(?:\n|$)`)
	match := re.FindStringSubmatch(text)

	if len(match) > 1 {
		// Bersihkan nilai yang didapatkan
		value := strings.TrimSpace(match[1])
		form["Kode unik"] = value
	}

	// Periksa apakah field wajib terisi
	if form["Kode unik"] == "" {
		return false, nil
	}

	return true, form
}

// processUpload memproses unggahan bukti transaksi
func (c *UploadProofCommand) processUpload(msg *message.Message, transactionCode string) (string, error) {
	// Validasi format kode transaksi (k_xxx00_000 atau m_xxx00_000)
	validCodePattern := regexp.MustCompile(`^[km]_[a-z]{3}\d{2}_\d{3}$`)
	if !validCodePattern.MatchString(transactionCode) {
		return fmt.Sprintf("❌ Format kode transaksi '%s' tidak valid. Format yang benar: k_mmm00_000 atau m_mmm00_000", transactionCode), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Unduh media
	mediaPath, err := msg.DownloadMedia()
	if err != nil {
		return fmt.Sprintf("❌ Gagal mengunduh bukti transaksi: %v", err), nil
	}

	// Pastikan file akan dihapus setelah selesai
	defer os.Remove(mediaPath)

	// Unggah bukti transaksi
	record, err := c.financeService.UploadTransactionProof(ctx, transactionCode, mediaPath)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			return fmt.Sprintf("❌ Transaksi dengan kode %s tidak ditemukan.", transactionCode), nil
		}
		if strings.Contains(err.Error(), "sudah memiliki bukti") {
			return fmt.Sprintf("❌ Transaksi dengan kode %s sudah memiliki bukti transaksi.", transactionCode), nil
		}
		return fmt.Sprintf("❌ Gagal mengunggah bukti transaksi: %v", err), nil
	}

	// Format response sukses menggunakan format baru yang lebih user-friendly
	recordDTO := dto.FromFinanceRecord(record)
	recordType := "PENGELUARAN"
	if record.Type == finance.TypeIncome {
		recordType = "PEMASUKAN"
	}

	// Tambahkan field berdasarkan tipe transaksi
	paymentMethodText := ""
	if record.Type == finance.TypeExpense {
		paymentMethodText = fmt.Sprintf("💳 Metode: %s\n", record.PaymentMethod)
	}

	// Gunakan nama lengkap dari field DTO
	storageTypeText := "Media Penyimpanan"
	if record.Type == finance.TypeExpense {
		storageTypeText = "Sumber Dana"
	}

	result := fmt.Sprintf(`────────────────────────
✅ BUKTI TRANSAKSI BERHASIL DITAMBAHKAN ✅
────────────────────────
Hai, pengguna 👋!
Bukti transaksi Anda berhasil diunggah.
📌 DETAIL %s: 
────────────────────────
📅 Tanggal: %s
📖 Deskripsi: %s
💰 Nominal: Rp %s
🏷 Kategori: %s
%s🏦 %s: %s
📝 Catatan: %s
📄 Bukti Transaksi: ✅ Tersedia
🔗 %s
────────────────────────`,
		recordType,
		recordDTO.DateFormatted,
		record.Description,
		recordDTO.AmountText,
		record.Category,
		paymentMethodText,
		storageTypeText,
		record.StorageMedia,
		record.Notes,
		record.ProofURL)

	return result, nil
}
