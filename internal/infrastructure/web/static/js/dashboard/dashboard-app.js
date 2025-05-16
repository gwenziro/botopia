/**
 * Dashboard Application
 * 
 * Mengelola data dan logika untuk dashboard Botopia
 */

document.addEventListener('alpine:init', () => {
    Alpine.data('dashboardApp', () => ({
        isConnected: false,
        connectedPhone: '',
        stats: {
            messageCount: 0,
            commandsRun: 0,
            commandCount: 0,
            uptime: 0
        },
        commands: {},
        loadingCommands: true,
        pollingInterval: null,

        initialize() {
            console.log('Initializing dashboard app');
            // Ambil data statistik dan status koneksi
            this.fetchStats();
            
            // Parse commands data dari element script
            this.fetchCommands();
            
            // Setup polling untuk memperbarui data secara berkala
            this.pollingInterval = setInterval(() => {
                this.fetchStats();
            }, 10000); // Update every 10 seconds

            // Bersihkan interval saat komponen dihapus
            this.$cleanup = () => {
                clearInterval(this.pollingInterval);
            };
        },

        fetchStats() {
            fetch('/api/stats')
                .then(response => response.json())
                .then(data => {
                    console.log('Stats loaded:', data);
                    this.isConnected = data.isConnected;
                    this.connectedPhone = data.phone || '';
                    this.stats = {
                        messageCount: data.messageCount || 0,
                        commandsRun: data.commandsRun || 0,
                        commandCount: data.commandCount || 0,
                        uptime: data.uptime || 0
                    };
                })
                .catch(error => {
                    console.error('Error fetching stats:', error);
                    if (typeof showToast === 'function') {
                        showToast('error', 'Gagal memuat data statistik');
                    }
                });
        },

        fetchCommands() {
            this.loadingCommands = true;
            console.log('Fetching commands data');
            
            try {
                // Get commands from embedded JSON data
                const commandsDataElement = document.getElementById('commands-data');
                if (commandsDataElement && commandsDataElement.textContent) {
                    const jsonText = commandsDataElement.textContent.trim();
                    console.log('Raw commands JSON snippet:', jsonText.substring(0, 100));
                    
                    try {
                        const commandsData = JSON.parse(jsonText);
                        console.log('Commands parsed successfully');
                        
                        if (commandsData && typeof commandsData === 'object' && Object.keys(commandsData).length > 0) {
                            this.commands = commandsData;
                            console.log('Commands loaded, count:', Object.keys(this.commands).length);
                            console.log('First command:', Object.keys(this.commands)[0]);
                        } else {
                            console.warn('Commands data is empty or invalid, falling back to API');
                            this.fetchCommandsFromApi();
                        }
                    } catch (parseError) {
                        console.error('Failed to parse commands JSON:', parseError);
                        this.fetchCommandsFromApi();
                    }
                } else {
                    console.warn('Commands data element not found or empty');
                    this.fetchCommandsFromApi();
                }
            } catch (error) {
                console.error('Error in fetchCommands:', error);
                this.commands = {};
            } finally {
                this.loadingCommands = false;
            }
        },

        fetchCommandsFromApi() {
            console.log('Fetching commands from API');
            // Fallback option: fetch commands from API
            fetch('/api/commands')
                .then(response => response.json())
                .then(data => {
                    console.log('Commands loaded from API:', Object.keys(data).length);
                    this.commands = data;
                })
                .catch(error => {
                    console.error('Error fetching commands from API:', error);
                    this.commands = {};
                })
                .finally(() => {
                    this.loadingCommands = false;
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
        },

        // Update teks tombol refresh berdasarkan status koneksi
        get refreshButtonText() {
            return this.isConnected ? 'Perbarui Data' : 'Perbarui Status';
        }
    }));
});
