package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gwenziro/botopia/internal/domain/contact"
	"github.com/gwenziro/botopia/internal/domain/repository"
	"github.com/gwenziro/botopia/internal/domain/service"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
)

// ContactService implementasi layanan kontak
type ContactService struct {
	contactRepo repository.ContactRepository
	log         *logger.Logger
}

// NewContactService membuat instance layanan kontak baru
func NewContactService(contactRepo repository.ContactRepository, log *logger.Logger) service.ContactService {
	return &ContactService{
		contactRepo: contactRepo,
		log:         log,
	}
}

// GetContact mendapatkan kontak berdasarkan nomor telepon
func (s *ContactService) GetContact(ctx context.Context, phone string) (*contact.Contact, error) {
	return s.contactRepo.FindByPhone(ctx, phone)
}

// GetAllContacts mendapatkan semua kontak
func (s *ContactService) GetAllContacts(ctx context.Context) ([]*contact.Contact, error) {
	return s.contactRepo.FindAll(ctx)
}

// AddContact menambahkan kontak baru
func (s *ContactService) AddContact(ctx context.Context, phone, name, notes string) (*contact.Contact, error) {
	// Validasi nomor telepon
	if phone == "" {
		return nil, fmt.Errorf("nomor telepon tidak boleh kosong")
	}

	// Normalisasi nomor telepon
	phone = normalizePhone(phone)

	// Cek apakah kontak sudah ada
	existingContact, err := s.contactRepo.FindByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}

	if existingContact != nil {
		return nil, fmt.Errorf("kontak dengan nomor %s sudah ada", phone)
	}

	// Buat kontak baru
	newContact := contact.NewContact(phone, name, notes)

	// Simpan kontak
	err = s.contactRepo.Save(ctx, newContact)
	if err != nil {
		return nil, err
	}

	return newContact, nil
}

// UpdateContact memperbarui kontak
func (s *ContactService) UpdateContact(ctx context.Context, phone, name, notes string, isActive bool) (*contact.Contact, error) {
	// Normalisasi nomor telepon
	phone = normalizePhone(phone)

	// Cek apakah kontak ada
	existingContact, err := s.contactRepo.FindByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}

	if existingContact == nil {
		return nil, fmt.Errorf("kontak dengan nomor %s tidak ditemukan", phone)
	}

	// Update data kontak
	existingContact.Name = name
	existingContact.Notes = notes
	existingContact.IsActive = isActive
	existingContact.UpdatedAt = time.Now()

	// Simpan perubahan
	err = s.contactRepo.Save(ctx, existingContact)
	if err != nil {
		return nil, err
	}

	return existingContact, nil
}

// DeleteContact menghapus kontak
func (s *ContactService) DeleteContact(ctx context.Context, phone string) error {
	// Normalisasi nomor telepon
	phone = normalizePhone(phone)
	return s.contactRepo.Delete(ctx, phone)
}

// IsWhitelisted memeriksa apakah nomor telepon diizinkan
func (s *ContactService) IsWhitelisted(ctx context.Context, phone string) (bool, error) {
	// Normalisasi nomor telepon
	phone = normalizePhone(phone)
	return s.contactRepo.IsWhitelisted(ctx, phone)
}

// GetWhitelistedContacts mendapatkan semua kontak dalam whitelist
func (s *ContactService) GetWhitelistedContacts(ctx context.Context) ([]*contact.Contact, error) {
	return s.contactRepo.GetWhitelistedContacts(ctx)
}

// SetWhitelistStatus mengatur status whitelist kontak
func (s *ContactService) SetWhitelistStatus(ctx context.Context, phone string, isActive bool) error {
	// Normalisasi nomor telepon
	phone = normalizePhone(phone)
	return s.contactRepo.SetWhitelistStatus(ctx, phone, isActive)
}

// normalizePhone menormalisasi format nomor telepon
func normalizePhone(phone string) string {
	// Implementasi sederhana, bisa diganti dengan library khusus
	if len(phone) > 0 && phone[0] == '+' {
		return phone
	}

	// Tambahkan + jika belum ada
	return "+" + strings.TrimPrefix(phone, "0")
}
