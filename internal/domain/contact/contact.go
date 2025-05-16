package contact

import "time"

// Contact merepresentasikan entitas kontak WhatsApp
type Contact struct {
	ID        string    `json:"id"`
	Phone     string    `json:"phone"`
	Name      string    `json:"name"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	IsActive  bool      `json:"isActive"` // Digunakan untuk whitelist
}

// NewContact membuat instance kontak baru
func NewContact(phone, name, notes string) *Contact {
	now := time.Now()
	return &Contact{
		Phone:     phone,
		Name:      name,
		Notes:     notes,
		CreatedAt: now,
		UpdatedAt: now,
		IsActive:  true,
	}
}
