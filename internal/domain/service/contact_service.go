package service

import (
	"context"

	"github.com/gwenziro/botopia/internal/domain/contact"
)

// ContactService mendefinisikan layanan untuk manajemen kontak
type ContactService interface {
	// GetContact mendapatkan kontak berdasarkan nomor telepon
	GetContact(ctx context.Context, phone string) (*contact.Contact, error)

	// GetAllContacts mendapatkan semua kontak
	GetAllContacts(ctx context.Context) ([]*contact.Contact, error)

	// AddContact menambahkan kontak baru
	AddContact(ctx context.Context, phone, name, notes string) (*contact.Contact, error)

	// UpdateContact memperbarui kontak
	UpdateContact(ctx context.Context, phone, name, notes string, isActive bool) (*contact.Contact, error)

	// DeleteContact menghapus kontak
	DeleteContact(ctx context.Context, phone string) error

	// IsWhitelisted memeriksa apakah nomor telepon diizinkan
	IsWhitelisted(ctx context.Context, phone string) (bool, error)

	// GetWhitelistedContacts mendapatkan semua kontak dalam whitelist
	GetWhitelistedContacts(ctx context.Context) ([]*contact.Contact, error)

	// SetWhitelistStatus mengatur status whitelist kontak
	SetWhitelistStatus(ctx context.Context, phone string, isActive bool) error
}
