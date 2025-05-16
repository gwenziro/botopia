/**
 * QR Scanner Module for Botopia
 * 
 * Handles QR code generation, display and connection management
 */

document.addEventListener('alpine:init', () => {
  Alpine.data('qrHandler', () => ({
    qrCode: '',
    isConnected: false,
    isRefreshing: false,
    connectionStatus: 'Terputus',
    statusTitle: 'Koneksi WhatsApp',
    statusMessage: 'Menunggu koneksi...',
    connectedPhone: '',
    deviceName: '',
    platform: '',
    businessName: '',
    deviceID: '',
    connectedSince: '',
    formattedConnectedTime: '',
    errorMessage: '',
    pollInterval: null,
    
    // Computed properties for UI
    get connectionClass() {
      if (this.isConnected) return 'connected';
      if (this.isRefreshing) return 'loading';
      return 'disconnected';
    },
    
    get connectionIcon() {
      if (this.isConnected) return 'fa-check';
      if (this.isRefreshing) return 'fa-spinner fa-spin';
      return 'fa-times';
    },
    
    get connectionBadgeClass() {
      if (this.isConnected) return 'connected';
      if (this.isRefreshing) return 'loading';
      return 'disconnected';
    },
    
    get displayName() {
      if (this.businessName) {
        return this.businessName;
      }
      return this.deviceName || 'WhatsApp User';
    },

    get deviceDescription() {
      let parts = [];
      if (this.platform) parts.push(this.platform);
      if (this.deviceID) parts.push(`ID: ${this.deviceID}`);
      return parts.join(' • ');
    },
    
    // Initialize QR handling
    initQR() {
      console.log('Initializing QR handler');
      this.fetchQRStatus();
      
      // Set up polling
      this.pollInterval = setInterval(() => {
        this.fetchQRStatus();
      }, 5000);
      
      // Clean up on page leave
      window.addEventListener('beforeunload', () => {
        if (this.pollInterval) {
          clearInterval(this.pollInterval);
        }
      });
    },
    
    // Fetch QR code and connection status
    fetchQRStatus() {
      fetch('/api/qr')
        .then(response => {
          if (!response.ok) {
            throw new Error('Network response was not ok');
          }
          return response.json();
        })
        .then(data => {
          console.log('QR status fetched:', data);
          
          // Update connection state
          this.isConnected = data.connectionState;
          
          // Update connection info
          if (this.isConnected) {
            this.connectionStatus = 'Terhubung';
            this.statusTitle = 'WhatsApp Terhubung';
            this.connectedPhone = data.phone || '';
            
            // Olah data perangkat dengan lebih baik
            this.deviceName = data.name || 'WhatsApp User';
            
            // Determine best platform display name
            if (data.deviceInfo && data.deviceInfo.platform) {
              this.platform = data.deviceInfo.platform;
            } else if (data.platform) {
              this.platform = data.platform;
            } else {
              // Fallback platform info with client detection
              const isMobile = /iPhone|iPad|iPod|Android/i.test(navigator.userAgent);
              const isTablet = /iPad|Android(?!.*Mobile)/i.test(navigator.userAgent);
              
              if (isMobile && !isTablet) {
                this.platform = 'WhatsApp Mobile';
              } else if (isTablet) {
                this.platform = 'WhatsApp Tablet';
              } else {
                this.platform = 'WhatsApp Web';
              }
            }
            
            this.businessName = data.businessName || '';
            this.deviceID = data.deviceID || '';

            // Format the connected time if available
            if (data.connectedSince) {
              this.connectedSince = data.connectedSince;
              this.formattedConnectedTime = this.formatConnectedTime(data.connectedSince);
            } else {
              // Jika tidak ada connectedSince dari server, gunakan waktu lokal
              if (!this.formattedConnectedTime) {
                this.formattedConnectedTime = this.formatTimestamp();
              }
            }

            // Create message with phone number if available
            if (this.connectedPhone) {
              this.statusMessage = `Terhubung ke ${this.connectedPhone}`;
            } else {
              this.statusMessage = 'WhatsApp terhubung';
            }
          } else {
            this.connectionStatus = 'Terputus';
            this.statusTitle = 'WhatsApp Terputus';
            this.statusMessage = 'Scan QR code untuk menghubungkan';
            this.connectedPhone = '';
            this.deviceName = '';
            this.platform = '';
            this.businessName = '';
            this.deviceID = '';
            this.connectedSince = '';
            this.formattedConnectedTime = '';
          }
          
          // Handle QR code
          if (data.qrCode && !this.isConnected) {
            if (this.qrCode !== data.qrCode) {
              this.qrCode = data.qrCode;
              this.renderQRCode();
            }
          } else {
            this.qrCode = '';
          }
          
          this.isRefreshing = false;
          this.errorMessage = '';
        })
        .catch(error => {
          console.error('Error fetching QR status:', error);
          this.errorMessage = 'Gagal memuat QR code: ' + error.message;
          this.isRefreshing = false;
        });
    },
    
    // Manually request new QR code
    refreshQR() {
      if (this.isRefreshing) return;
      
      this.isRefreshing = true;
      this.qrCode = '';
      this.statusMessage = 'Memperbarui QR code...';
      
      // Force refresh with cache-busting parameter
      fetch('/api/qr?refresh=' + new Date().getTime())
        .then(response => {
          if (!response.ok) {
            throw new Error('Failed to refresh QR code');
          }
          return response.json();
        })
        .then(data => {
          if (data.qrCode) {
            this.qrCode = data.qrCode;
            this.renderQRCode();
            this.statusMessage = 'QR code berhasil diperbarui';
            
            // Show success toast
            showToast('success', 'QR code berhasil diperbarui');
          } else {
            throw new Error('No QR code received');
          }
        })
        .catch(error => {
          console.error('Error refreshing QR:', error);
          this.errorMessage = 'Gagal memperbarui QR code: ' + error.message;
          this.statusMessage = 'Gagal memperbarui QR code';
          
          // Show error toast
          showToast('error', 'Gagal memperbarui QR code');
        })
        .finally(() => {
          this.isRefreshing = false;
        });
    },
    
    // Disconnect WhatsApp session
    disconnectWhatsApp() {
      if (!confirm('Yakin ingin memutuskan koneksi WhatsApp?')) {
        return;
      }
      
      fetch('/api/disconnect', {
        method: 'POST'
      })
        .then(response => {
          if (!response.ok) {
            throw new Error('Failed to disconnect');
          }
          return response.json();
        })
        .then(data => {
          this.isConnected = false;
          this.connectionStatus = 'Terputus';
          this.statusTitle = 'WhatsApp Terputus';
          this.statusMessage = 'Koneksi terputus';
          this.connectedPhone = '';
          this.qrCode = '';
          this.deviceName = '';
          this.platform = '';
          this.businessName = '';
          this.deviceID = '';
          this.connectedSince = '';
          this.formattedConnectedTime = '';
          
          // Show success toast
          showToast('success', 'WhatsApp berhasil diputuskan');
          
          // Request new QR code after short delay
          setTimeout(() => {
            this.refreshQR();
          }, 1500);
        })
        .catch(error => {
          console.error('Error disconnecting:', error);
          showToast('error', 'Gagal memutuskan koneksi: ' + error.message);
        });
    },
    
    // Format the connected time from ISO string
    formatConnectedTime(isoString) {
      try {
        const date = new Date(isoString);
        return date.toLocaleString('id-ID', {
          day: 'numeric',
          month: 'long',
          year: 'numeric',
          hour: '2-digit',
          minute: '2-digit',
        });
      } catch (e) {
        console.error('Error formatting date:', e);
        return this.formatTimestamp();
      }
    },
    
    // Render QR code using qrcode.js library
    renderQRCode() {
      if (!this.qrCode) return;
      
      const container = document.getElementById('qrcode');
      if (!container) return;
      
      // Clear previous QR code
      container.innerHTML = '';
      
      // Create new QR code
      new QRCode(container, {
        text: this.qrCode,
        width: 250,
        height: 250,
        colorDark: '#000000',
        colorLight: '#ffffff',
        correctLevel: QRCode.CorrectLevel.H
      });
      
      console.log('QR code rendered');
    },
    
    // Helper for timestamp formatting
    formatTimestamp() {
      const now = new Date();
      return now.toLocaleString('id-ID', { 
        day: 'numeric',
        month: 'long',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    }
  }));
});
