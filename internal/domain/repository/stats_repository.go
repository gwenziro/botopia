package repository

// StatsRepository mendefinisikan kontrak untuk statistik bot
type StatsRepository interface {
	// IncrementMessageCount menambah hitungan pesan
	IncrementMessageCount() error

	// IncrementCommandCount menambah hitungan command
	IncrementCommandCount() error

	// GetStats mendapatkan statistik
	GetStats() (*BotStats, error)

	// SetConnectionState mengatur status koneksi
	SetConnectionState(isConnected bool) error
}

// BotStats berisi statistik bot
type BotStats struct {
	ConnectionState string
	IsConnected     bool
	MessageCount    int
	CommandsRun     int
	Uptime          int64
	SystemUptime    int64
	ConnectedSince  int64
}
