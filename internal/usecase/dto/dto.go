package dto

// StatsDTO adalah data transfer object untuk statistik bot
type StatsDTO struct {
	ConnectionState string            `json:"connectionState"`
	IsConnected     bool              `json:"isConnected"`
	MessageCount    int               `json:"messageCount"`
	CommandsRun     int               `json:"commandsRun"`
	Uptime          int64             `json:"uptime"`
	SystemUptime    int64             `json:"systemUptime"`
	CommandCount    int               `json:"commandCount"`
	Phone           string            `json:"phone,omitempty"`
	Name            string            `json:"name,omitempty"`
	PushName        string            `json:"pushName,omitempty"`
	DeviceDetails   *DeviceDetailsDTO `json:"deviceDetails,omitempty"`
}

// DeviceDetailsDTO mewakili detail perangkat
type DeviceDetailsDTO struct {
	Platform    string `json:"platform,omitempty"`
	DeviceModel string `json:"deviceModel,omitempty"`
	OSVersion   string `json:"osVersion,omitempty"`
	ClientType  string `json:"clientType,omitempty"`
	IPAddress   string `json:"ipAddress,omitempty"`
	DeviceID    string `json:"deviceId,omitempty"`
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
