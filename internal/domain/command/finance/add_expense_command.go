package finance

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gwenziro/botopia/internal/domain/dto"
	"github.com/gwenziro/botopia/internal/domain/finance"
	"github.com/gwenziro/botopia/internal/domain/message"
	"github.com/gwenziro/botopia/internal/domain/service"
	"github.com/gwenziro/botopia/internal/utils"
)

// AddExpenseCommand implementasi command untuk menambahkan pengeluaran
type AddExpenseCommand struct {
	financeService service.FinanceService
}

// NewAddExpenseCommand membuat instance command baru
func NewAddExpenseCommand(financeService service.FinanceService) *AddExpenseCommand {
	return &AddExpenseCommand{
		financeService: financeService,
	}
}

// GetName mengembalikan nama command
func (c *AddExpenseCommand) GetName() string {
	return "keluar"
}

// GetDescription mengembalikan deskripsi command
func (c *AddExpenseCommand) GetDescription() string {
	return "Mencatat pengeluaran baru. Kirim !keluar untuk mendapatkan form input data."
}

// Execute menjalankan command
func (c *AddExpenseCommand) Execute(args []string, msg *message.Message) (string, error) {
	// Jika tidak ada argumen, kirimkan form template
	if len(args) == 0 {
		return c.getFormTemplate(), nil
	}

	// Cek apakah pesan adalah form yang diisi
	if isFilledForm, form := c.parseFormInput(msg.Text); isFilledForm {
		// Periksa apakah ada media yang dilampirkan untuk upload bukti
		var mediaPath string
		var err error

		if msg.HasMedia() {
			// Download media jika ada
			mediaPath, err = msg.DownloadMedia()
			if err != nil {
				return fmt.Sprintf("Gagal mengunduh media: %v", err), nil
			}
			// Pastikan file akan dihapus setelah selesai
			defer func() {
				if mediaPath != "" {
					os.Remove(mediaPath)
				}
			}()
		}

		return c.processForm(form, mediaPath)
	}

	// Jika bukan form dan ada argument, tampilkan panduan
	config, _ := c.financeService.GetConfiguration(context.Background())
	helpMsg := "Untuk mencatat pengeluaran, kirim !keluar (tanpa parameter) untuk mendapatkan formulir."

	if config != nil {
		// Tambahkan informasi kategori yang tersedia
		helpMsg += "\n\nKategori tersedia: " + strings.Join(config.ExpenseCategories, ", ")
		helpMsg += "\nMetode pembayaran: " + strings.Join(config.PaymentMethods, ", ")
		helpMsg += "\nSumber dana: " + strings.Join(config.StorageMedias, ", ")
	}

	return helpMsg, nil
}

// getFormTemplate mengembalikan template form pengeluaran
func (c *AddExpenseCommand) getFormTemplate() string {
	return `!keluar
────────────────────────
💰 INPUT DATA PENGELUARAN 💰
────────────────────────
Tanggal: 
Deskripsi: 
Nominal: 
Kategori: 
Metode: 
Sumber: 
Catatan: 
────────────────────────
Tolong catat data pengeluaranku di atas, ya! 🙏`
}

// parseFormInput mengecek dan mem-parsing input formulir
func (c *AddExpenseCommand) parseFormInput(text string) (bool, map[string]string) {
	// Cek apakah dimulai dengan !keluar
	if !strings.HasPrefix(text, "!keluar") {
		return false, nil
	}

	// Regex untuk mengekstrak nilai dari setiap field
	form := make(map[string]string)

	// Ekstrak nilai field dari text dengan regex yang lebih presisi
	fields := []string{"Tanggal", "Deskripsi", "Nominal", "Kategori", "Metode", "Sumber", "Catatan"}
	for i, field := range fields {
		var pattern string
		if i == len(fields)-1 {
			// Pola khusus untuk field terakhir (Catatan) - hentikan sebelum separator
			pattern = fmt.Sprintf(`%s:\s*(.*?)(?:\n-+|$)`, regexp.QuoteMeta(field))
		} else {
			// Pola umum untuk field lainnya
			pattern = fmt.Sprintf(`%s:\s*(.*?)(?:\n%s:|$)`, regexp.QuoteMeta(field), regexp.QuoteMeta(fields[i+1]))
		}

		re := regexp.MustCompile(pattern)
		match := re.FindStringSubmatch(text)

		if len(match) > 1 {
			// Bersihkan nilai yang didapatkan
			value := strings.TrimSpace(match[1])
			form[field] = value
		}
	}

	// Periksa apakah minimal field wajib terisi
	requiredFields := []string{"Tanggal", "Deskripsi", "Nominal", "Kategori", "Metode", "Sumber"}
	for _, field := range requiredFields {
		if form[field] == "" {
			return false, nil
		}
	}

	return true, form
}

// processForm memproses form yang sudah diisi
func (c *AddExpenseCommand) processForm(form map[string]string, mediaPath string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse tanggal
	date, err := utils.ParseDateWithFormats(form["Tanggal"])
	if err != nil {
		return fmt.Sprintf("Format tanggal tidak valid: %v. Gunakan format seperti '15 Mei 2025'", err), nil
	}

	// Parse nominal
	amount, err := utils.ParseMoney(form["Nominal"])
	if err != nil {
		return fmt.Sprintf("Nominal tidak valid: %v. Gunakan angka saja, contoh: 50000", err), nil
	}

	description := form["Deskripsi"]
	category := form["Kategori"]
	paymentMethod := form["Metode"]
	storageMedia := form["Sumber"]

	// Tangani catatan kosong dengan satu strip
	notes := form["Catatan"]
	if notes == "" || strings.Contains(notes, "─") {
		notes = "-"
	}

	// Validasi parameter terhadap konfigurasi
	if err := c.financeService.ValidateAddExpenseParams(ctx, category, paymentMethod, storageMedia); err != nil {
		return fmt.Sprintf("Validasi gagal: %v", err), nil
	}

	var record *finance.FinanceRecord

	// Jika ada media path, unggah terlebih dahulu sebelum menyimpan record
	if mediaPath != "" {
		// Buat record dengan URL kosong terlebih dahulu
		tmpRecord, err := c.financeService.AddExpenseWithDate(
			ctx, date, description, amount, category,
			paymentMethod, storageMedia, notes, "",
		)

		if err != nil {
			return fmt.Sprintf("Gagal mencatat pengeluaran: %v", err), nil
		}

		// Unggah bukti menggunakan kode transaksi yang dihasilkan
		uploadCtx, uploadCancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer uploadCancel()

		record, err = c.financeService.UploadTransactionProof(uploadCtx, tmpRecord.UniqueCode, mediaPath)
		if err != nil {
			// Transaksi sudah tersimpan tapi gagal upload bukti
			proofStatus := fmt.Sprintf("\n\n⚠️ Gagal mengunggah bukti: %v", err)
			return c.formatSuccessResponse(tmpRecord, false) + proofStatus, nil
		}
	} else {
		// Tanpa media, langsung simpan record
		record, err = c.financeService.AddExpenseWithDate(
			ctx, date, description, amount, category,
			paymentMethod, storageMedia, notes, "",
		)

		if err != nil {
			return fmt.Sprintf("Gagal mencatat pengeluaran: %v", err), nil
		}
	}

	// Format response sukses
	return c.formatSuccessResponse(record, mediaPath != ""), nil
}

// formatSuccessResponse memformat pesan sukses
func (c *AddExpenseCommand) formatSuccessResponse(record *finance.FinanceRecord, hasProof bool) string {
	recordDTO := dto.FromFinanceRecord(record)
	proofStatus := "Belum tersedia"
	if hasProof || recordDTO.HasProof {
		proofStatus = "✅ Tersedia"
	}

	result := fmt.Sprintf(`────────────────────────
✅ DATA PENGELUARAN BERHASIL DITAMBAHKAN ✅
────────────────────────
Data pengeluaran kamu berhasil dicatat.

📌 DETAIL PENGELUARAN: 
────────────────────────
📅 Tanggal: %s
📖 Deskripsi: %s
💰 Jumlah: Rp %s
🏷 Kategori: %s
💳 Metode: %s
🏦 Sumber Dana: %s
📝 Catatan: %s
🧾 Bukti Transaksi: %s
────────────────────────
ℹ Kode Transaksi: %s
Gunakan kode ini untuk melampirkan bukti transaksi di kemudian hari.
────────────────────────
💡 Ketik !ringkasan untuk melihat laporan keuanganmu! 📊
────────────────────────`,
		recordDTO.DateFormatted,
		record.Description,
		recordDTO.AmountText,
		record.Category,
		record.PaymentMethod,
		record.StorageMedia,
		record.Notes,
		proofStatus,
		record.UniqueCode)

	return result
}
