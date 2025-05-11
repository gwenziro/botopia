/**
 * QR Code page functionality
 * Digunakan oleh halaman QR untuk menampilkan dan memperbarui kode QR
 */

// Define the Alpine component function globally
window.qrPage = function() {
    return {
        qrCode: '',
        connectionState: 'disconnected',
        loading: false,
        refreshing: false,
        pollingInterval: null,
        lastQrCode: '',
        qrRetryCount: 0,
        phone: '', // Phone number of connected account
        name: '',  // Name of connected account
        
        init() {
            console.log("QR Page initialized");
            
            // Get initial values from data attributes
            const qrApp = document.getElementById('qr-app');
            if (qrApp) {
                this.connectionState = qrApp.dataset.connectionState || 'disconnected';
                this.qrCode = qrApp.dataset.qrCode || '';
                this.lastQrCode = this.qrCode;
                this.phone = qrApp.dataset.phone || '';
                this.name = qrApp.dataset.name || '';
                
                console.log("Initial QR data loaded:", { 
                    hasQR: !!this.qrCode, 
                    state: this.connectionState,
                    phone: this.phone,
                    name: this.name
                });
            }
            
            // Start polling for QR code updates
            this.startPolling();
        },
        
        startPolling() {
            // Fetch immediately on first load
            this.fetchQRCode();
            
            // Then poll every 10 seconds
            this.pollingInterval = setInterval(() => {
                this.fetchQRCode();
            }, 10000);
            
            console.log("QR polling started");
        },
        
        stopPolling() {
            if (this.pollingInterval) {
                clearInterval(this.pollingInterval);
                this.pollingInterval = null;
                console.log("QR polling stopped");
            }
        },
        
        fetchQRCode() {
            // Don't fetch if we're already refreshing
            if (this.refreshing) return;
            
            this.loading = true;
            this.refreshing = true;
            
            fetch('/api/qr')
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response error');
                    }
                    return response.json();
                })
                .then(data => {
                    // Track connection state changes
                    const wasConnected = this.connectionState === 'connected';
                    const isNowConnected = data.connectionState === true;
                    
                    console.log("QR API response:", { 
                        hasQRCode: !!data.qrCode, 
                        connectionState: data.connectionState,
                        phone: data.phone,
                        name: data.name
                    });
                    
                    // Only update QR code if we receive a non-empty one
                    if (data.qrCode) {
                        // Check if this is a new QR code
                        const isNewQr = this.lastQrCode !== data.qrCode;
                        this.qrCode = data.qrCode;
                        
                        // If it's a new QR code, show notification and reset retry count
                        if (isNewQr) {
                            this.lastQrCode = data.qrCode;
                            this.qrRetryCount = 0;
                            
                            window.showToast({
                                title: 'New QR Code',
                                message: 'A new QR code is available for scanning',
                                type: 'info',
                                duration: 5000
                            });
                        }
                    } else {
                        // If no QR code but we're disconnected, increment retry count
                        if (!isNowConnected && this.qrRetryCount < 3) {
                            this.qrRetryCount++;
                            console.log(`No QR code received, retry ${this.qrRetryCount}/3`);
                        }
                    }
                    
                    // Update contact information if available
                    if (data.phone) this.phone = data.phone;
                    if (data.name) this.name = data.name;
                    
                    // Update connection state
                    this.connectionState = isNowConnected ? 'connected' : 'disconnected';
                    
                    // If connection state changed, show notification
                    if (wasConnected !== isNowConnected) {
                        window.showToast({
                            title: isNowConnected ? 'Connected!' : 'Disconnected',
                            message: isNowConnected 
                                ? 'Successfully connected to WhatsApp'
                                : 'Connection to WhatsApp has been lost',
                            type: isNowConnected ? 'success' : 'error',
                            duration: 5000
                        });
                    }
                })
                .catch(error => {
                    console.error('Error fetching QR code:', error);
                    window.showToast({
                        title: 'Error',
                        message: 'Failed to fetch QR code from server',
                        type: 'error'
                    });
                })
                .finally(() => {
                    this.loading = false;
                    setTimeout(() => {
                        this.refreshing = false;
                    }, 1000);
                });
        },
        
        refreshQR() {
            if (this.refreshing) return;
            this.fetchQRCode();
            
            // Show a notification that we're refreshing
            window.showToast({
                title: 'Refreshing QR Code',
                message: 'Fetching a fresh QR code from the server',
                type: 'info',
                duration: 3000
            });
        },
        
        getQRImageUrl() {
            if (!this.qrCode) return '';
            return `https://api.qrserver.com/v1/create-qr-code/?data=${encodeURIComponent(this.qrCode)}&size=300x300`;
        },

        // Ensure a safe display for phone value
        getDisplayPhone() {
            if (!this.phone || this.phone === "") {
                return "Unknown number";
            }
            return this.phone;
        },

        // Tetap tampilkan nama default
        getDisplayName() {
            return "WhatsApp User";
        }
    };
};

// Make sure the component gets initialized when the page loads
document.addEventListener('DOMContentLoaded', function() {
    // Add debugging message to help diagnose issues
    console.log("QR Scanner script loaded");
});
