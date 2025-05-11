package dto

// StatsDTO adalah data transfer object untuk statistik bot
type StatsDTO struct {
	ConnectionState string `json:"connectionState"`
	IsConnected     bool   `json:"isConnected"`
	MessageCount    int    `json:"messageCount"`
	CommandsRun     int    `json:"commandsRun"`
	Uptime          int64  `json:"uptime"`
	CommandCount    int    `json:"commandCount"`
	Phone           string `json:"phone,omitempty"`
	Name            string `json:"name,omitempty"`
}

// ConnectionStatusDTO adalah data transfer object untuk status koneksi
type ConnectionStatusDTO struct {
	IsConnected bool   `json:"isConnected"`
	Message     string `json:"message"`
	Phone       string `json:"phone,omitempty"`
	Error       error  `json:"-"` // Tidak diserialisasi
}

// QRCodeDTO adalah data transfer object untuk QR code
type QRCodeDTO struct {
	QRCode          string `json:"qrCode"`
	ConnectionState bool   `json:"connectionState"`
	Phone           string `json:"phone,omitempty"`
	Name            string `json:"name,omitempty"`
}

// CommandDTO adalah data transfer object untuk command
type CommandDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Usage       string `json:"usage,omitempty"`
	Category    string `json:"category,omitempty"`
}
