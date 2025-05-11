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

	// Hitung uptime jika terhubung
	if stats.ConnectionState == "connected" {
		stats.Uptime = int64(time.Since(time.Unix(stats.ConnectedSince, 0)).Seconds())
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
		r.stats.ConnectedSince = time.Now().Unix()
	} else {
		r.stats.ConnectionState = "disconnected"
	}

	return nil
}
