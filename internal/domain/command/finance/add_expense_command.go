package finance

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gwenziro/botopia/internal/domain/message"
	"github.com/gwenziro/botopia/internal/domain/service"
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
		return c.processForm(form)
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
func (c *AddExpenseCommand) processForm(form map[string]string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse tanggal
	date, err := parseFormDate(form["Tanggal"])
	if err != nil {
		return fmt.Sprintf("Format tanggal tidak valid: %v. Gunakan format seperti '15 Mei 2025'", err), nil
	}

	// Parse nominal
	amount, err := parseFormAmount(form["Nominal"])
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

	// Tambahkan pengeluaran dengan tanggal custom
	record, err := c.financeService.AddExpenseWithDate(
		ctx,
		date,
		description,
		amount,
		category,
		paymentMethod,
		storageMedia,
		notes, // Pastikan notes sudah diproses dengan benar
		"",    // ProofURL kosong
	)

	if err != nil {
		return fmt.Sprintf("Gagal mencatat pengeluaran: %v", err), nil
	}

	// Format response sukses dengan urutan parameter yang benar
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
⚠ Bukti Transaksi: Belum tersedia
────────────────────────
ℹ Kode Transaksi: %s
Gunakan kode ini untuk melampirkan bukti transaksi di kemudian hari.
────────────────────────
💡 Ketik !ringkasan untuk melihat laporan keuanganmu! 📊
────────────────────────`,
		formatDateOutput(date),
		description,
		formatMoney(amount),
		category,
		paymentMethod,
		storageMedia,
		notes,
		record.UniqueCode)

	return result, nil
}

// Helper functions
func parseFormDate(dateStr string) (time.Time, error) {
	// Support berbagai format tanggal umum
	formats := []string{
		"2 Jan 2006",
		"2 January 2006",
		"02 Jan 2006",
		"02 January 2006",
		"2006-01-02",
		"02/01/2006",
		"2/1/2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	// Format bulan dalam bahasa Indonesia
	indoMonths := map[string]string{
		"jan": "Jan", "feb": "Feb", "mar": "Mar", "apr": "Apr",
		"mei": "May", "jun": "Jun", "jul": "Jul", "agu": "Aug",
		"sep": "Sep", "okt": "Oct", "nov": "Nov", "des": "Dec",
		"januari": "January", "februari": "February", "maret": "March", "april": "April",
		"juni": "June", "juli": "July", "agustus": "August",
		"september": "September", "oktober": "October", "november": "November", "desember": "December",
	}

	// Coba parse format Indonesia (misal: "15 Mei 2025")
	re := regexp.MustCompile(`(\d{1,2})\s+([A-Za-z]+)\s+(\d{4})`)
	match := re.FindStringSubmatch(dateStr)
	if len(match) == 4 {
		day, month, year := match[1], strings.ToLower(match[2]), match[3]
		if englishMonth, ok := indoMonths[month]; ok {
			newDateStr := fmt.Sprintf("%s %s %s", day, englishMonth, year)
			for _, format := range []string{"2 Jan 2006", "2 January 2006"} {
				if t, err := time.Parse(format, newDateStr); err == nil {
					return t, nil
				}
			}
		}
	}

	// Bila tanggal kosong atau tidak valid, gunakan tanggal hari ini
	if dateStr == "" || dateStr == "hari ini" {
		return time.Now(), nil
	}

	return time.Time{}, fmt.Errorf("format tanggal tidak dikenali")
}

func parseFormAmount(amountStr string) (float64, error) {
	// Bersihkan string dari karakter non-numerik kecuali titik dan koma
	numericStr := regexp.MustCompile(`[^0-9.,]`).ReplaceAllString(amountStr, "")

	// Ganti koma dengan titik untuk format float
	numericStr = strings.Replace(numericStr, ",", ".", -1)

	return strconv.ParseFloat(numericStr, 64)
}

func formatDateOutput(date time.Time) string {
	// Format: "07 Mei 2025"
	indoMonths := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}

	return fmt.Sprintf("%02d %s %d", date.Day(), indoMonths[date.Month()-1], date.Year())
}

func formatMoney(amount float64) string {
	// Format: "50.000" (dengan pemisah ribuan)
	str := strconv.FormatFloat(amount, 'f', 0, 64)
	result := ""

	for i, c := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += "."
		}
		result += string(c)
	}

	return result
}
