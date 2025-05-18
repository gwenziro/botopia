/**
 * Connectivity Application
 * Mengelola antarmuka untuk pemindai QR WhatsApp dan status koneksi
 */
document.addEventListener('alpine:init', () => {
    Alpine.data('connectivityApp', () => ({
        connectionState: 'disconnected', // connected, connecting, disconnected
        qrCode: '',
        phone: '',
        name: 'WhatsApp User',
        uptime: 0,
        connectedSince: null,
        autoRefresh: true,
        pollingInterval: null,
        retryCount: 0,
        maxRetries: 15,
        deviceDetails: null,
        isLoading: false,
        
        initialize() {
            console.log('Initializing Connectivity app');
            this.loadInitialData();
            this.startPolling();
            
            // Cleanup when component is destroyed
            this.$once('$destroy', () => {
                this.stopPolling();
            });
        },
        
        loadInitialData() {
            // Dapatkan data dari Alpine props jika tersedia
            const alpineProps = this.$el.dataset;
            
            // Load connection state
            if (alpineProps.connectionState) {
                this.connectionState = alpineProps.connectionState === 'true' ? 'connected' : 'disconnected';
            }
            
            if (alpineProps.qrCode) {
                this.qrCode = alpineProps.qrCode;
            }
            
            if (alpineProps.phone) {
                this.phone = alpineProps.phone;
            }
            
            if (alpineProps.name) {
                this.name = alpineProps.name;
            }
            
            // Load uptime jika ada
            if (alpineProps.uptime) {
                this.uptime = parseInt(alpineProps.uptime, 10);
                if (this.uptime > 0) {
                    this.connectedSince = new Date(Date.now() - (this.uptime * 1000));
                }
            }
            
            // Segera fetch data untuk informasi terbaru
            this.fetchConnectivityData();
        },
        
        startPolling() {
            // Jika sudah connected, polling lebih jarang
            const interval = this.connectionState === 'connected' ? 10000 : 3000;
            
            // Clear any existing interval
            this.stopPolling();
            
            // Start new polling
            this.pollingInterval = setInterval(() => {
                this.fetchConnectivityData();
            }, interval);
        },
        
        stopPolling() {
            if (this.pollingInterval) {
                clearInterval(this.pollingInterval);
                this.pollingInterval = null;
            }
        },
        
        fetchConnectivityData() {
            if (this.isLoading) return;
            
            this.isLoading = true;
            
            fetch('/api/qr')
                .then(response => {
                    if (!response.ok) {
                        throw new Error(`HTTP error: ${response.status}`);
                    }
                    return response.json();
                })
                .then(data => {
                    // Update state based on connection status
                    const wasConnected = this.connectionState === 'connected';
                    this.connectionState = data.connectionState ? 'connected' : 'disconnected';
                    
                    // If connection state changed to connected, show notification
                    if (!wasConnected && this.connectionState === 'connected') {
                        this.showNotification('success', 'Berhasil terhubung ke WhatsApp!');
                    } else if (wasConnected && this.connectionState !== 'connected') {
                        this.showNotification('error', 'Koneksi WhatsApp terputus!');
                    }
                    
                    // If connected, update device info
                    if (this.connectionState === 'connected') {
                        this.phone = data.phone || '';
                        this.name = data.name || 'WhatsApp User';
                        this.uptime = data.uptime || 0;
                        this.deviceDetails = data.deviceDetails || null;
                        
                        // Calculate connected since from uptime
                        if (this.uptime > 0) {
                            this.connectedSince = new Date(Date.now() - (this.uptime * 1000));
                        }
                        
                        // Reset retry counter
                        this.retryCount = 0;
                        
                        // Slow down polling when connected
                        this.adjustPollingInterval(10000);
                    } else {
                        // If not connected, try to get QR code
                        if (data.qrCode) {
                            this.qrCode = data.qrCode;
                            // Reset retry counter when QR code received
                            this.retryCount = 0;
                            
                            // Speed up polling when waiting for scan
                            this.adjustPollingInterval(3000);
                        } else {
                            // Increment retry counter
                            this.retryCount++;
                            
                            // Stop polling if exceeded max retries
                            if (this.retryCount > this.maxRetries) {
                                this.stopPolling();
                                this.showNotification('error', 'Gagal mendapatkan QR code, silakan refresh halaman');
                            }
                        }
                    }
                    
                    this.isLoading = false;
                })
                .catch(error => {
                    console.error('Error fetching connectivity data:', error);
                    this.retryCount++;
                    
                    if (this.retryCount > this.maxRetries) {
                        this.stopPolling();
                        this.showNotification('error', 'Terjadi error saat mengambil data. Silakan refresh halaman');
                    }
                    
                    this.isLoading = false;
                });
        },
        
        showNotification(type, message) {
            // Use global toast function if available
            if (typeof showToast === 'function') {
                showToast(type, message);
            } else {
                // Fallback to alert
                alert(message);
            }
        },
        
        adjustPollingInterval(newInterval) {
            this.stopPolling();
            this.pollingInterval = setInterval(() => {
                this.fetchConnectivityData();
            }, newInterval);
        },
        
        disconnect() {
            if (confirm('Apakah Anda yakin ingin memutuskan koneksi WhatsApp?')) {
                fetch('/api/disconnect', {
                    method: 'POST',
                })
                .then(response => {
                    if (!response.ok) {
                        throw new Error(`HTTP error: ${response.status}`);
                    }
                    return response.json();
                })
                .then(data => {
                    if (data.success) {
                        this.connectionState = 'disconnected';
                        this.qrCode = '';
                        this.deviceDetails = null;
                        this.showNotification('success', 'Koneksi WhatsApp berhasil diputuskan');
                        
                        // Restart polling with faster interval
                        this.adjustPollingInterval(3000);
                    } else {
                        this.showNotification('error', data.error || 'Gagal memutuskan koneksi');
                    }
                })
                .catch(error => {
                    console.error('Error disconnecting:', error);
                    this.showNotification('error', 'Terjadi kesalahan saat memutuskan koneksi');
                });
            }
        },
        
        formatUptime(seconds) {
            if (!seconds || seconds <= 0) {
                return '0 menit';
            }
            
            const days = Math.floor(seconds / 86400);
            const hours = Math.floor((seconds % 86400) / 3600);
            const minutes = Math.floor((seconds % 3600) / 60);
            
            let result = [];
            if (days > 0) result.push(`${days} hari`);
            if (hours > 0) result.push(`${hours} jam`);
            if (minutes > 0 || (!days && !hours)) result.push(`${minutes} menit`);
            
            return result.join(' ');
        },
        
        formatDate(date) {
            if (!date) return '-';
            
            try {
                // Format tanggal dalam bahasa Indonesia: "25 Agustus 2023, 14:30"
                const options = { 
                    day: 'numeric', 
                    month: 'long', 
                    year: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit'
                };
                
                return new Date(date).toLocaleDateString('id-ID', options);
            } catch (error) {
                console.error('Error formatting date:', error);
                return new Date(date).toString();
            }
        },
        
        formatPhone(phone) {
            if (!phone) return '-';
            // Format as international number if not already formatted
            if (!phone.startsWith('+')) {
                return '+' + phone;
            }
            return phone;
        },
        
        getQRImageURL() {
            if (!this.qrCode) return '';
            const encodedQR = encodeURIComponent(this.qrCode);
            return `https://chart.googleapis.com/chart?chs=300x300&cht=qr&chl=${encodedQR}&choe=UTF-8`;
        }
    }));
});
