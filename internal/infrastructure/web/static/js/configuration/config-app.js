/**
 * Configuration Application
 * 
 * Mengelola data dan logika untuk halaman konfigurasi Botopia
 */

document.addEventListener('alpine:init', () => {
    Alpine.data('configApp', () => ({
        isConnected: false,
        isGoogleConnected: false,
        connectedPhone: '',
        spreadsheetUrl: 'https://docs.google.com/spreadsheets/d/',
        stats: {
            uptime: 0
        },
        config: {
            commandPrefix: '!',
            logLevel: 'INFO',
            webPort: 8080,
            webHost: '0.0.0.0',
            spreadsheetId: '',
            driveFolderId: '',
            credentialsFile: './service-account.json'
        },

        initialize() {
            // Ambil data statistik dan status koneksi
            this.fetchStats();
            
            // Ambil konfigurasi saat ini
            this.fetchConfig();
        },

        fetchStats() {
            fetch('/api/stats')
                .then(response => response.json())
                .then(data => {
                    this.isConnected = data.isConnected;
                    this.connectedPhone = data.phone || '';
                    this.stats = {
                        uptime: data.uptime || 0
                    };
                })
                .catch(error => {
                    console.error('Error fetching stats:', error);
                });
        },

        fetchConfig() {
            fetch('/api/config')
                .then(response => response.json())
                .then(data => {
                    // Update local config with server data
                    this.config = {
                        ...this.config,
                        ...data
                    };
                    
                    // Set additional status fields
                    this.isGoogleConnected = !!data.spreadsheetId;
                    
                    if (data.spreadsheetId) {
                        this.spreadsheetUrl = `https://docs.google.com/spreadsheets/d/${data.spreadsheetId}`;
                    }
                })
                .catch(error => {
                    console.error('Error fetching config:', error);
                    
                    // Simulate success for demo
                    this.isGoogleConnected = true;
                });
        },

        saveConfig() {
            // In a real application, we'd post the config to the server
            fetch('/api/config', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(this.config),
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                showToast('success', 'Konfigurasi berhasil disimpan');
            })
            .catch(error => {
                console.error('Error saving config:', error);
                showToast('error', 'Gagal menyimpan konfigurasi');
            });
        },

        formatUptime(seconds) {
            if (!seconds || seconds <= 0) return 'Tidak aktif';
            
            // Konversi detik ke format yang lebih ramah
            const days = Math.floor(seconds / 86400);
            const hours = Math.floor((seconds % 86400) / 3600);
            const minutes = Math.floor((seconds % 3600) / 60);
            
            if (days > 0) {
                return `${days}d ${hours}h ${minutes}m`;
            } else if (hours > 0) {
                return `${hours}h ${minutes}m`;
            } else if (minutes > 0) {
                return `${minutes}m`;
            } else {
                return `${seconds}s`;
            }
        }
    }));
});
