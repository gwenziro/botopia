/**
 * Dashboard Application
 * Menangani fungsionalitas dashboard utama
 */
document.addEventListener('alpine:init', () => {
    Alpine.data('dashboardApp', () => ({
        stats: {
            messageCount: 0,
            commandsRun: 0,
            commandCount: 0,
            systemUptime: 0
        },
        isConnected: false,
        phone: '',
        name: 'WhatsApp User',
        uptime: 0,
        connectedSince: null,
        deviceDetails: null,
        commands: {},
        
        // Command to Icon Mapping
        commandIcons: {
            ping: 'fa-heartbeat',
            panduan: 'fa-question-circle',
            help: 'fa-info-circle',
            keluar: 'fa-money-bill-wave',
            masuk: 'fa-wallet',
            unggah: 'fa-upload',
            ringkasan: 'fa-chart-pie'
        },
        
        // Google Service Properties - Enhanced with IDs
        googleServices: {
            isConfigured: false,
            sheets: {
                isConnected: false,
                url: '#',
                id: ''
            },
            drive: {
                isConnected: false,
                url: '#',
                id: ''
            }
        },
        
        initialize() {
            console.log('Initializing dashboard app');
            
            // Load data awal
            this.loadInitialData();
            
            // Set up polling untuk update otomatis
            this.startPolling();
            
            // Cleanup when component is destroyed
            window.addEventListener('beforeunload', () => this.stopPolling());
        },
        
        loadInitialData() {
            // Dapatkan data dari Alpine props jika tersedia
            this.loadPropsData();
            
            // Muat stats via API untuk memastikan data terbaru
            this.fetchStats();
            
            // Load commands
            this.fetchCommands();
            
            // Load Google service status
            this.fetchGoogleServicesStatus();
        },
        
        loadPropsData() {
            const alpineProps = this.$el.dataset;
            
            // Google services configuration
            if (alpineProps.hasGoogleApiService !== undefined) {
                this.googleServices.isConfigured = alpineProps.hasGoogleApiService === 'true';
            }
            
            if (alpineProps.hasFinanceService !== undefined) {
                this.googleServices.sheets.isConnected = alpineProps.hasFinanceService === 'true';
                this.googleServices.drive.isConnected = alpineProps.hasFinanceService === 'true';
            }
            
            if (alpineProps.spreadsheetUrl) {
                this.googleServices.sheets.url = alpineProps.spreadsheetUrl;
            }
            
            if (alpineProps.spreadsheetId) {
                this.googleServices.sheets.id = alpineProps.spreadsheetId;
            }
            
            if (alpineProps.driveFolderId) {
                this.googleServices.drive.id = alpineProps.driveFolderId;
                this.googleServices.drive.url = `https://drive.google.com/drive/folders/${alpineProps.driveFolderId}`;
            }
            
            // Command data from JSON string
            if (alpineProps.commandsJson) {
                try {
                    this.commands = JSON.parse(alpineProps.commandsJson);
                } catch (e) {
                    console.error('Error parsing commands JSON', e);
                    this.commands = {};
                }
            }
            
            // Connection status
            if (alpineProps.isConnected) {
                this.isConnected = alpineProps.isConnected === 'true';
            }
            
            // User details
            if (alpineProps.phone) {
                this.phone = alpineProps.phone;
            }
            
            if (alpineProps.name) {
                this.name = alpineProps.name;
            }
            
            // Stats
            if (alpineProps.messageCount) {
                this.stats.messageCount = parseInt(alpineProps.messageCount, 10) || 0;
            }
            
            if (alpineProps.commandsRun) {
                this.stats.commandsRun = parseInt(alpineProps.commandsRun, 10) || 0;
            }
            
            if (alpineProps.commandCount) {
                this.stats.commandCount = parseInt(alpineProps.commandCount, 10) || 0;
            }
            
            if (alpineProps.uptime) {
                this.uptime = parseInt(alpineProps.uptime, 10) || 0;
                if (this.uptime > 0) {
                    this.connectedSince = new Date(Date.now() - (this.uptime * 1000));
                }
            }
        },
        
        startPolling() {
            // Poll for stats update every 30 seconds
            this.statsInterval = setInterval(() => this.fetchStats(), 30000);
            
            // Poll for service status less frequently
            this.serviceInterval = setInterval(() => this.fetchGoogleServicesStatus(), 60000);
        },
        
        stopPolling() {
            if (this.statsInterval) clearInterval(this.statsInterval);
            if (this.serviceInterval) clearInterval(this.serviceInterval);
        },
        
        fetchStats() {
            fetch('/api/stats')
                .then(response => response.json())
                .then(data => {
                    // Update stats values
                    this.stats.messageCount = data.messageCount || 0;
                    this.stats.commandsRun = data.commandsRun || 0;
                    this.stats.commandCount = data.commandCount || 0;
                    this.stats.systemUptime = data.systemUptime || 0;
                    
                    // Update connection info
                    this.uptime = data.uptime || 0;
                    this.isConnected = data.isConnected || false;
                    this.phone = data.phone || '';
                    
                    // Update device details
                    this.deviceDetails = data.deviceDetails || null;
                    
                    // Update name dengan benar
                    if (data.name && data.name !== 'WhatsApp User') {
                        this.name = data.name;
                    } else if (data.pushName && data.pushName !== '') {
                        this.name = data.pushName;
                    } else {
                        this.name = data.phone || 'WhatsApp User';
                    }
                    
                    // Calculate connected since from uptime
                    if (this.uptime > 0) {
                        this.connectedSince = new Date(Date.now() - (this.uptime * 1000));
                    }
                })
                .catch(error => {
                    console.error('Error fetching stats:', error);
                });
        },
        
        fetchGoogleServicesStatus() {
            fetch('/api/config/status')
                .then(response => response.json())
                .then(data => {
                    console.log("Google services status:", data);
                    
                    // Update Google services status
                    this.googleServices.isConfigured = data.googleApi?.configured || false;
                    this.googleServices.sheets.isConnected = data.googleApi?.sheets || false;
                    this.googleServices.drive.isConnected = data.googleApi?.drive || false;
                    
                    // Store full IDs (not truncated)
                    if (data.spreadsheetId) {
                        this.googleServices.sheets.id = data.spreadsheetId;
                    }
                    
                    if (data.spreadsheetUrl) {
                        this.googleServices.sheets.url = data.spreadsheetUrl;
                    }
                    
                    if (data.driveFolderId) {
                        this.googleServices.drive.id = data.driveFolderId;
                    }
                    
                    if (data.driveFolderUrl) {
                        this.googleServices.drive.url = data.driveFolderUrl;
                    } else if (data.driveFolderId) {
                        this.googleServices.drive.url = `https://drive.google.com/drive/folders/${data.driveFolderId}`;
                    }
                })
                .catch(error => {
                    console.error('Error fetching Google services status:', error);
                });
        },
        
        fetchCommands() {
            // If commands are already loaded, don't fetch again
            if (Object.keys(this.commands).length > 0) return;
            
            fetch('/api/commands')
                .then(response => response.json())
                .then(data => {
                    this.commands = data;
                })
                .catch(error => {
                    console.error('Error fetching commands:', error);
                });
        },
        
        // Get command icon based on name
        getCommandIcon(name) {
            return this.commandIcons[name] || 'fa-code'; // Default to code icon
        },
        
        // Format ID for display (truncate if too long)
        formatId(id) {
            if (!id) return '-';
            if (id.length > 12) {
                return id.substring(0, 6) + '...' + id.substring(id.length - 6);
            }
            return id;
        },
        
        formatPhone(phone) {
            if (!phone) return '-';
            // Format as international number if not already formatted
            if (!phone.startsWith('+')) {
                return '+' + phone;
            }
            return phone;
        },
        
        formatNumber(num) {
            return new Intl.NumberFormat().format(num);
        },
        
        formatDate(date) {
            if (!date) return '-';
            
            try {
                return new Date(date).toLocaleString('id-ID', {
                    day: 'numeric',
                    month: 'long',
                    year: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit'
                });
            } catch (e) {
                console.error('Error formatting date:', e);
                return '-';
            }
        },
        
        formatUptimeDuration(seconds) {
            if (!seconds || seconds <= 0) return '0 menit';
            
            const days = Math.floor(seconds / 86400);
            const hours = Math.floor((seconds % 86400) / 3600);
            const minutes = Math.floor((seconds % 3600) / 60);
            
            let result = [];
            if (days > 0) result.push(`${days} hari`);
            if (hours > 0) result.push(`${hours} jam`);
            if (minutes > 0 || (!days && !hours)) result.push(`${minutes} menit`);
            
            return result.join(' ');
        }
    }));
});
