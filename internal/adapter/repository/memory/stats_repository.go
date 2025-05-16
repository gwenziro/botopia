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
			ConnectionState: "disconnected",
			MessageCount:    0,
			CommandsRun:     0,
			ConnectedSince:  0,
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

	// Buat salinan untuk menghindari race condition
	stats := r.stats

	// Hitung uptime sistem (berbeda dari uptime koneksi)
	stats.SystemUptime = int64(time.Since(r.startTime).Seconds())

	// Hitung uptime jika terhubung
	if stats.ConnectionState == "connected" {
		if stats.ConnectedSince > 0 {
			stats.Uptime = int64(time.Since(time.Unix(stats.ConnectedSince, 0)).Seconds())
		} else {
			// Fallback ke system uptime jika ConnectedSince tidak diset tapi state connected
			stats.Uptime = stats.SystemUptime
			// Setel ConnectedSince ke waktu sekarang dikurangi SystemUptime
			stats.ConnectedSince = time.Now().Unix() - stats.SystemUptime
		}
	} else {
		stats.Uptime = 0
	}

	return &stats, nil
}

// SetConnectionState mengatur status koneksi
func (r *StatsRepository) SetConnectionState(isConnected bool) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if isConnected {
		r.stats.ConnectionState = "connected"
		// Hanya set ConnectedSince jika sebelumnya tidak connected
		if r.stats.ConnectedSince == 0 {
			r.stats.ConnectedSince = time.Now().Unix()
		}
	} else {
		r.stats.ConnectionState = "disconnected"
		r.stats.ConnectedSince = 0
		r.stats.Uptime = 0
	}

	return nil
}
