package user

// User merepresentasikan pengguna WhatsApp
type User struct {
	ID            string
	Name          string
	Phone         string
	PushName      string
	IsBot         bool
	DeviceDetails *DeviceDetails
}

// DeviceDetails berisi informasi detail perangkat yang terhubung
type DeviceDetails struct {
	Platform     string // OS/Platform perangkat
	BusinessName string // Nama bisnis jika akun bisnis
	DeviceID     string // ID perangkat
	Connected    string // Waktu terhubung dalam format RFC3339
	DeviceModel  string // Model perangkat (baru)
	OSVersion    string // Versi OS (baru)
	ClientType   string // Tipe client (baru)
	IPAddress    string // Alamat IP (baru)
}
