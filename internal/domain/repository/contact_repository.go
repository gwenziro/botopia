package repository

import (
	"context"

	"github.com/gwenziro/botopia/internal/domain/contact"
)

// ContactRepository mendefinisikan kontrak untuk repository kontak
type ContactRepository interface {
	// FindByPhone mencari kontak berdasarkan nomor telepon
	FindByPhone(ctx context.Context, phone string) (*contact.Contact, error)

	// FindAll mendapatkan semua kontak
	FindAll(ctx context.Context) ([]*contact.Contact, error)

	// Save menyimpan kontak baru atau memperbarui yang sudah ada
	Save(ctx context.Context, contact *contact.Contact) error

	// Delete menghapus kontak
	Delete(ctx context.Context, phone string) error

	// IsWhitelisted memeriksa apakah nomor telepon ada dalam whitelist
	IsWhitelisted(ctx context.Context, phone string) (bool, error)

	// GetWhitelistedContacts mendapatkan semua kontak dalam whitelist
	GetWhitelistedContacts(ctx context.Context) ([]*contact.Contact, error)

	// SetWhitelistStatus mengatur status whitelist kontak
	SetWhitelistStatus(ctx context.Context, phone string, isActive bool) error
}
