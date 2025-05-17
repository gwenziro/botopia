package file

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gwenziro/botopia/internal/domain/contact"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
)

// ContactRepository implementasi repository kontak yang menyimpan data di file JSON
type ContactRepository struct {
	contacts map[string]*contact.Contact // In-memory cache
	mutex    sync.RWMutex
	filePath string
	log      *logger.Logger
}

// NewContactRepository membuat instance repository kontak baru
func NewContactRepository(dataDir string, log *logger.Logger) *ContactRepository {
	filePath := filepath.Join(dataDir, "contacts.json")
	repo := &ContactRepository{
		contacts: make(map[string]*contact.Contact),
		filePath: filePath,
		log:      log,
	}

	// Load data dari file saat inisialisasi
	repo.loadContacts()

	return repo
}

// loadContacts memuat data kontak dari file
func (r *ContactRepository) loadContacts() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Cek apakah file ada
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		r.log.Info("File kontak tidak ditemukan: %s, membuat baru", r.filePath)
		return
	}

	// Baca file
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		r.log.Error("Gagal membaca file kontak: %v", err)
		return
	}

	// Unmarshal JSON
	var contacts []*contact.Contact
	if err := json.Unmarshal(data, &contacts); err != nil {
		r.log.Error("Gagal parse data kontak: %v", err)
		return
	}

	// Simpan ke map
	for _, c := range contacts {
		r.contacts[c.Phone] = c
	}

	r.log.Info("Berhasil memuat %d kontak dari file", len(r.contacts))
}

// saveContacts menyimpan data kontak ke file
func (r *ContactRepository) saveContacts() error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Konversi map ke slice untuk JSON
	var contacts []*contact.Contact
	for _, c := range r.contacts {
		contacts = append(contacts, c)
	}

	// Marshal ke JSON
	data, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		return err
	}

	// Buat direktori jika belum ada
	dir := filepath.Dir(r.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Tulis ke file
	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return err
	}

	r.log.Info("Berhasil menyimpan %d kontak ke file", len(r.contacts))
	return nil
}

// FindByPhone mencari kontak berdasarkan nomor telepon
func (r *ContactRepository) FindByPhone(ctx context.Context, phone string) (*contact.Contact, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if c, exists := r.contacts[phone]; exists {
		return c, nil
	}

	return nil, nil // Tidak ditemukan, bukan error
}

// FindAll mendapatkan semua kontak
func (r *ContactRepository) FindAll(ctx context.Context) ([]*contact.Contact, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	contacts := make([]*contact.Contact, 0, len(r.contacts))
	for _, c := range r.contacts {
		contacts = append(contacts, c)
	}

	return contacts, nil
}

// Save menyimpan kontak baru atau memperbarui yang sudah ada
func (r *ContactRepository) Save(ctx context.Context, c *contact.Contact) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Update waktu update
	c.UpdatedAt = time.Now()

	r.contacts[c.Phone] = c
	r.log.Info("Kontak disimpan: %s (%s)", c.Name, c.Phone)

	// Simpan ke file setiap kali ada perubahan
	go func() {
		if err := r.saveContacts(); err != nil {
			r.log.Error("Gagal menyimpan kontak ke file: %v", err)
		}
	}()

	return nil
}

// Delete menghapus kontak
func (r *ContactRepository) Delete(ctx context.Context, phone string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.contacts[phone]; exists {
		delete(r.contacts, phone)
		r.log.Info("Kontak dihapus: %s", phone)

		// Simpan perubahan ke file
		go func() {
			if err := r.saveContacts(); err != nil {
				r.log.Error("Gagal menyimpan kontak ke file setelah penghapusan: %v", err)
			}
		}()

		return nil
	}

	return nil // Tidak ada kontak yang dihapus, bukan error
}

// IsWhitelisted memeriksa apakah nomor telepon ada dalam whitelist
func (r *ContactRepository) IsWhitelisted(ctx context.Context, phone string) (bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if c, exists := r.contacts[phone]; exists {
		return c.IsActive, nil
	}

	return false, nil
}

// GetWhitelistedContacts mendapatkan semua kontak dalam whitelist
func (r *ContactRepository) GetWhitelistedContacts(ctx context.Context) ([]*contact.Contact, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	contacts := make([]*contact.Contact, 0)
	for _, c := range r.contacts {
		if c.IsActive {
			contacts = append(contacts, c)
		}
	}

	return contacts, nil
}

// SetWhitelistStatus mengatur status whitelist kontak
func (r *ContactRepository) SetWhitelistStatus(ctx context.Context, phone string, isActive bool) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if c, exists := r.contacts[phone]; exists {
		c.IsActive = isActive
		c.UpdatedAt = time.Now()
		r.log.Info("Status whitelist kontak diperbarui: %s (%s) = %v", c.Name, phone, isActive)

		// Simpan perubahan ke file
		go func() {
			if err := r.saveContacts(); err != nil {
				r.log.Error("Gagal menyimpan kontak ke file setelah update whitelist: %v", err)
			}
		}()

		return nil
	}

	// Jika tidak ada, buat baru
	newContact := contact.NewContact(phone, phone, "")
	newContact.IsActive = isActive
	r.contacts[phone] = newContact
	r.log.Info("Kontak baru ditambahkan ke whitelist: %s", phone)

	// Simpan perubahan ke file
	go func() {
		if err := r.saveContacts(); err != nil {
			r.log.Error("Gagal menyimpan kontak ke file setelah tambah whitelist: %v", err)
		}
	}()

	return nil
}
