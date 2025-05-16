package memory

import (
	"context"
	"sync"
	"time"

	"github.com/gwenziro/botopia/internal/domain/contact"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
)

// ContactRepository implementasi in-memory untuk repository kontak
type ContactRepository struct {
	contacts map[string]*contact.Contact
	mutex    sync.RWMutex
	log      *logger.Logger
}

// NewContactRepository membuat instance repository kontak baru
func NewContactRepository(log *logger.Logger) *ContactRepository {
	return &ContactRepository{
		contacts: make(map[string]*contact.Contact),
		log:      log,
	}
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

	return nil
}

// Delete menghapus kontak
func (r *ContactRepository) Delete(ctx context.Context, phone string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.contacts[phone]; exists {
		delete(r.contacts, phone)
		r.log.Info("Kontak dihapus: %s", phone)
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
		return nil
	}

	// Jika tidak ada, buat baru
	newContact := contact.NewContact(phone, phone, "")
	newContact.IsActive = isActive
	r.contacts[phone] = newContact
	r.log.Info("Kontak baru ditambahkan ke whitelist: %s", phone)

	return nil
}
