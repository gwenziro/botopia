package finance

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gwenziro/botopia/internal/domain/message"
	"github.com/gwenziro/botopia/internal/domain/service"
)

// UploadProofCommand implementasi command untuk mengunggah bukti transaksi
type UploadProofCommand struct {
	financeService service.FinanceService
}

// NewUploadProofCommand membuat instance command baru
func NewUploadProofCommand(financeService service.FinanceService) *UploadProofCommand {
	return &UploadProofCommand{
		financeService: financeService,
	}
}

// GetName mengembalikan nama command
func (c *UploadProofCommand) GetName() string {
	return "unggah"
}

// GetDescription mengembalikan deskripsi command
func (c *UploadProofCommand) GetDescription() string {
	return "Mengunggah bukti transaksi untuk catatan yang sudah ada. Format: !unggah <kode_transaksi> (sertakan gambar/foto bukti)"
}

// Execute menjalankan command
func (c *UploadProofCommand) Execute(args []string, msg *message.Message) (string, error) {
	// Cek apakah ada gambar yang dilampirkan
	if !msg.HasMedia() {
		return "❌ Mohon lampirkan foto bukti transaksi untuk diunggah.", nil
	}

	// Cek apakah ada kode transaksi yang diberikan
	if len(args) == 0 {
		return "❌ Mohon sertakan kode transaksi. Format: !unggah <kode_transaksi>", nil
	}

	// Ambil kode transaksi
	transactionCode := args[0]

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

	// Format response sukses
	recordType := "pengeluaran"
	if record.Type == "income" {
		recordType = "pemasukan"
	}

	result := fmt.Sprintf(`✅ BUKTI TRANSAKSI BERHASIL DIUNGGAH ✅

Detail Transaksi:
📅 Tanggal: %s
📖 Deskripsi: %s
💰 Nominal: Rp %s
🏷️ Kategori: %s
🧾 Jenis: %s

🔗 Bukti transaksi telah tersimpan. Cek melalui tautan ini:
%s`,
		formatDateOutput(record.Date),
		record.Description,
		formatMoney(record.Amount),
		record.Category,
		recordType,
		record.ProofURL)

	return result, nil
}
