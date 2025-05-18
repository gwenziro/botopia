/**
 * Dashboard Application
 * 
 * Menampilkan overview dashboard Botopia
 */
document.addEventListener('alpine:init', () => {
    Alpine.data('dashboardApp', () => ({
        // Status data
        isConnected: false,
        connectionState: "disconnected",
        commandCount: 0,
        messageCount: 0,
        commandsRun: 0,
        phone: "",
        name: "WhatsApp User",
        uptime: 0,
        
        // Finance data
        transactions: [],
        weeklyStats: {
            totalIncome: 0,
            totalExpense: 0,
            balance: 0,
            largestCategory: "",
            largestAmount: 0,
            categoryData: {}
        },
        
        // Commands data
        commands: {},
        popularCommands: [],
        
        // Service flags
        hasFinanceService: false,
        hasGoogleAPIService: false,
        spreadsheetUrl: "#",
        
        // Chart instance
        expenseChart: null,
        
        initialize() {
            console.log('Initializing dashboard app...');
            
            // Load initial dashboard stats
            this.loadDashboardStats();
            
            // Set reload interval for stats (every 30 seconds)
            setInterval(() => this.loadDashboardStats(true), 30000);
            
            // Get finance data if available
            if (this.hasFinanceService) {
                this.loadTransactions();
                
                // Set service variables from HTML
                this.hasFinanceService = document.getElementById('finance-service-enabled')?.value === 'true';
                this.hasGoogleAPIService = document.getElementById('google-api-enabled')?.value === 'true';
                this.spreadsheetUrl = document.getElementById('spreadsheet-url')?.value || '#';
                
                // Initialize chart after a brief delay to ensure DOM is ready
                setTimeout(() => this.initializeChart(), 500);
            }
        },
        
        loadDashboardStats(isUpdate = false) {
            // Load from API endpoint
            fetch('/api/stats')
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Failed to fetch stats');
                    }
                    return response.json();
                })
                .then(data => {
                    this.updateStats(data);
                    this.loadCommands();
                    
                    // Show notification on updates
                    if (isUpdate && this.isConnected !== data.isConnected) {
                        if (data.isConnected) {
                            showToast('success', 'WhatsApp berhasil terhubung!');
                        } else {
                            showToast('warning', 'Koneksi WhatsApp terputus!');
                        }
                    }
                })
                .catch(error => {
                    console.error('Error loading dashboard stats:', error);
                });
        },
        
        updateStats(data) {
            this.isConnected = data.isConnected;
            this.connectionState = data.connectionState;
            this.commandCount = data.commandCount || 0;
            this.messageCount = data.messageCount || 0;
            this.commandsRun = data.commandsRun || 0;
            this.phone = data.phone || "";
            this.name = data.name || "WhatsApp User";
            this.uptime = data.uptime || 0;
            
            // Update title with connection status
            document.title = `Dashboard | ${this.isConnected ? '✓' : '✗'} Botopia`;
        },
        
        loadCommands() {
            // Ambil dari script tag JSON terlebih dahulu jika ada
            const cmdDataEl = document.getElementById('commands-data');
            if (cmdDataEl && cmdDataEl.textContent) {
                try {
                    const cmdData = JSON.parse(cmdDataEl.textContent);
                    this.processCommandData(cmdData);
                    return;
                } catch (e) {
                    console.warn('Failed to parse commands data from script tag:', e);
                }
            }
            
            // Fallback to API
            fetch('/api/commands')
                .then(response => response.json())
                .then(data => {
                    this.processCommandData(data);
                })
                .catch(error => {
                    console.error('Error loading commands:', error);
                });
        },
        
        processCommandData(data) {
            this.commands = data;
            
            // Extract popular commands (finance related and core commands)
            const commandList = [];
            for (const name in data) {
                commandList.push({
                    name: name,
                    description: data[name].description,
                    category: data[name].category || 'Uncategorized'
                });
            }
            
            // Sort by priority categories
            const priorityOrder = {
                'Keuangan': 1,
                'Finance': 1,
                'System': 2,
                'Help': 3
            };
            
            commandList.sort((a, b) => {
                const priorityA = priorityOrder[a.category] || 999;
                const priorityB = priorityOrder[b.category] || 999;
                return priorityA - priorityB;
            });
            
            // Take top 5
            this.popularCommands = commandList.slice(0, 5);
        },
        
        loadTransactions() {
            // Try to get from JSON embedded in page first
            const txDataEl = document.getElementById('transactions-data');
            if (txDataEl && txDataEl.textContent) {
                try {
                    const txData = JSON.parse(txDataEl.textContent);
                    this.transactions = txData;
                    
                    // Get weekly stats
                    const weeklyStatsEl = document.getElementById('weekly-stats-data');
                    if (weeklyStatsEl && weeklyStatsEl.textContent) {
                        this.weeklyStats = JSON.parse(weeklyStatsEl.textContent);
                    }
                    
                    return;
                } catch (e) {
                    console.warn('Failed to parse transactions data from script tag:', e);
                }
            }
            
            // Fallback to API
            fetch('/api/transactions/recent?limit=5')
                .then(response => response.json())
                .then(data => {
                    this.transactions = data.transactions || [];
                })
                .catch(error => {
                    console.error('Error loading transactions:', error);
                    this.transactions = [];
                });
        },
        
        initializeChart() {
            const chartContainer = document.getElementById('expense-chart');
            if (!chartContainer) return;
            
            // If Chart.js is not loaded, load it dynamically
            if (typeof Chart === 'undefined') {
                const script = document.createElement('script');
                script.src = 'https://cdn.jsdelivr.net/npm/chart.js';
                script.onload = () => this.createExpenseChart();
                document.head.appendChild(script);
            } else {
                this.createExpenseChart();
            }
        },
        
        createExpenseChart() {
            const categories = Object.keys(this.weeklyStats.categoryData || {});
            const amounts = categories.map(cat => this.weeklyStats.categoryData[cat]);
            
            if (categories.length === 0) {
                document.getElementById('expense-chart').innerHTML = 
                    '<div class="flex h-full items-center justify-center text-slate-500">Tidak ada data untuk ditampilkan</div>';
                return;
            }
            
            // Generate colors
            const colors = [
                'rgba(99, 102, 241, 0.7)',   // Indigo
                'rgba(245, 158, 11, 0.7)',   // Amber
                'rgba(16, 185, 129, 0.7)',   // Emerald
                'rgba(14, 165, 233, 0.7)',   // Sky
                'rgba(168, 85, 247, 0.7)',   // Purple
                'rgba(249, 115, 22, 0.7)',   // Orange
                'rgba(239, 68, 68, 0.7)',    // Red
                'rgba(236, 72, 153, 0.7)'    // Pink
            ];
            
            // Create chart
            const ctx = document.getElementById('expense-chart');
            this.expenseChart = new Chart(ctx, {
                type: 'bar',
                data: {
                    labels: categories,
                    datasets: [{
                        label: 'Pengeluaran',
                        data: amounts,
                        backgroundColor: colors.slice(0, categories.length),
                        borderColor: 'rgba(255, 255, 255, 0.2)',
                        borderWidth: 1
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    scales: {
                        y: {
                            beginAtZero: true,
                            grid: {
                                color: 'rgba(255, 255, 255, 0.1)'
                            },
                            ticks: {
                                color: 'rgba(255, 255, 255, 0.7)'
                            }
                        },
                        x: {
                            grid: {
                                display: false
                            },
                            ticks: {
                                color: 'rgba(255, 255, 255, 0.7)'
                            }
                        }
                    },
                    plugins: {
                        legend: {
                            display: false
                        }
                    }
                }
            });
        },
        
        // Formatter untuk uptime
        formatUptime(seconds) {
            if (!seconds) return "0 detik";
            
            const days = Math.floor(seconds / 86400);
            const hours = Math.floor((seconds % 86400) / 3600);
            const minutes = Math.floor((seconds % 3600) / 60);
            
            let result = [];
            if (days > 0) result.push(`${days} hari`);
            if (hours > 0) result.push(`${hours} jam`);
            if (minutes > 0) result.push(`${minutes} menit`);
            
            return result.join(", ");
        },
        
        // Format angka dengan pemisah ribuan
        formatNumber(num) {
            if (num === undefined || num === null) return "0";
            
            return new Intl.NumberFormat('id-ID').format(num);
        }
    }));
});
