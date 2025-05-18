package memory

import (
	"sync"
	"time"

	"github.com/gwenziro/botopia/internal/domain/repository"
)

// StatsRepository implementasi in-memory untuk statistik bot
type StatsRepository struct {
	stats     repository.BotStats
	mutex     sync.RWMutex
	startTime time.Time
}

// NewStatsRepository membuat instance stats repository baru
func NewStatsRepository() *StatsRepository {
	return &StatsRepository{
		stats: repository.BotStats{
			MessageCount: 0,
			CommandsRun:  0,
			IsConnected:  false,
		},
		startTime: time.Now(),
	}
}

// IncrementMessageCount menambah hitungan pesan
func (r *StatsRepository) IncrementMessageCount() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.stats.MessageCount++
	return nil
}

// IncrementCommandCount menambah hitungan command
func (r *StatsRepository) IncrementCommandCount() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.stats.CommandsRun++
	return nil
}

// GetStats mendapatkan statistik bot
func (r *StatsRepository) GetStats() (*repository.BotStats, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Salin data untuk thread safety
	stats := r.stats

	// Tambahkan system uptime
	systemUptime := time.Since(r.startTime).Seconds()
	stats.SystemUptime = int64(systemUptime)

	// Hitung uptime koneksi jika terhubung
	if stats.IsConnected && stats.ConnectedSince > 0 {
		uptime := time.Now().Unix() - stats.ConnectedSince
		stats.Uptime = uptime
	} else {
		stats.Uptime = 0
	}

	return &stats, nil
}

// SetConnectionState mengatur status koneksi
func (r *StatsRepository) SetConnectionState(isConnected bool) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Update state and connection details
	if isConnected && !r.stats.IsConnected {
		// Baru terhubung - catat waktu mulai
		r.stats.IsConnected = true
		r.stats.ConnectionState = "connected"
		r.stats.ConnectedSince = time.Now().Unix()
	} else if !isConnected && r.stats.IsConnected {
		// Koneksi terputus
		r.stats.IsConnected = false
		r.stats.ConnectionState = "disconnected"
		r.stats.ConnectedSince = 0
		r.stats.Uptime = 0
	}

	return nil
}
