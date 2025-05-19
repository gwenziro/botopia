/**
 * Connection Status Badge Component
 * Badge universal yang menampilkan status koneksi WhatsApp
 */
document.addEventListener('alpine:init', () => {
    Alpine.data('connectionStatusBadge', () => ({
        isConnected: false,
        checkInterval: null,
        isLoading: false,
        
        init() {
            console.log('Initializing connection status badge');
            this.checkConnectionStatus();
            
            // Set polling interval untuk update status otomatis
            this.checkInterval = setInterval(() => {
                this.checkConnectionStatus();
            }, 15000); // Check every 15 seconds
            
            // Listen for connection events from other components
            window.addEventListener('connection-status-changed', (event) => {
                if (event.detail && event.detail.connectionState !== undefined) {
                    this.isConnected = event.detail.connectionState;
                }
            });
            
            // Clean up interval when component is destroyed
            this.$cleanup = () => {
                if (this.checkInterval) {
                    clearInterval(this.checkInterval);
                    this.checkInterval = null;
                }
            };
        },
        
        checkConnectionStatus() {
            if (this.isLoading) return;
            
            this.isLoading = true;
            
            fetch('/api/qr')
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    return response.json();
                })
                .then(data => {
                    // Update connection status
                    const wasConnected = this.isConnected;
                    this.isConnected = data.connectionState === true;
                    
                    // If status changed, dispatch a global event
                    if (wasConnected !== this.isConnected) {
                        window.dispatchEvent(new CustomEvent('connection-status-changed', {
                            detail: { connectionState: this.isConnected }
                        }));
                    }
                    
                    this.isLoading = false;
                })
                .catch(error => {
                    console.error('Error checking connection status:', error);
                    this.isLoading = false;
                });
        }
    }));
});
