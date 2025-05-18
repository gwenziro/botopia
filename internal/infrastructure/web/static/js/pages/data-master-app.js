/**
 * Data Master Application - Read-Only Version
 * 
 * Menampilkan data master seperti kategori, metode pembayaran, dll.
 */
document.addEventListener('alpine:init', () => {
    Alpine.data('dataMasterApp', () => ({
        activeTab: 'expense-categories',
        loading: true,
        masterData: {
            expenseCategories: [],
            incomeCategories: [],
            paymentMethods: [],
            storageMedias: []
        },
        
        initializeDataMaster() {
            console.log('Initializing data master app');
            
            // Set active tab from URL if available
            const urlParams = new URLSearchParams(window.location.search);
            if (urlParams.has('tab')) {
                this.activeTab = urlParams.get('tab');
            }
            
            // Load data dari sumber terbaik yang tersedia
            this.loadMasterData();
            
            // Watch for tab changes to update URL
            this.$watch('activeTab', (value) => {
                const url = new URL(window.location);
                url.searchParams.set('tab', value);
                history.replaceState(null, '', url);
            });
        },
        
        // Fungsi refresh data - sama seperti loadMasterData
        refreshData() {
            // Jika sudah loading, abaikan
            if (this.loading) return;
            
            // Reset data dan load ulang dari API
            this.loading = true;
            this.fetchMasterDataFromAPI();
            
            // Tampilkan notifikasi
            showToast('info', 'Memperbarui data...');
        },
        
        loadMasterData() {
            this.loading = true;
            console.log('Loading master data...');
            
            // Strategi 1: Langsung gunakan data dari window.botopiaConfig jika tersedia
            if (window.botopiaConfig && window.botopiaConfig.masterData) {
                try {
                    console.log('Using data from global config object');
                    const md = window.botopiaConfig.masterData;
                    
                    this.masterData = {
                        expenseCategories: Array.isArray(md.expenseCategories) ? md.expenseCategories : [],
                        incomeCategories: Array.isArray(md.incomeCategories) ? md.incomeCategories : [],
                        paymentMethods: Array.isArray(md.paymentMethods) ? md.paymentMethods : [],
                        storageMedias: Array.isArray(md.storageMedias) ? md.storageMedias : []
                    };
                    
                    console.log('Master data loaded successfully from global object', this.masterData);
                    this.loading = false;
                    return;
                } catch (error) {
                    console.error('Error loading from global config:', error);
                }
            }
            
            // Strategi 2: Parse dari script tag JSON
            const configDataElement = document.getElementById('config-data');
            if (configDataElement && configDataElement.textContent.trim()) {
                try {
                    const jsonData = JSON.parse(configDataElement.textContent.trim());
                    console.log('Using data from script tag');
                    
                    this.masterData = {
                        expenseCategories: Array.isArray(jsonData.expenseCategories) ? jsonData.expenseCategories : [],
                        incomeCategories: Array.isArray(jsonData.incomeCategories) ? jsonData.incomeCategories : [],
                        paymentMethods: Array.isArray(jsonData.paymentMethods) ? jsonData.paymentMethods : [],
                        storageMedias: Array.isArray(jsonData.storageMedias) ? jsonData.storageMedias : []
                    };
                    
                    console.log('Master data loaded successfully from script tag', this.masterData);
                    this.loading = false;
                    return;
                } catch (error) {
                    console.error('Error parsing JSON from script tag:', error);
                }
            }
            
            // Strategi 3: Ambil dari API sebagai fallback
            console.log('Fetching data from API as fallback');
            this.fetchMasterDataFromAPI();
        },
        
        // Helper method to get array from object with multiple possible property names
        getArrayFromObject(obj, possibleNames) {
            if (!obj) return [];
            
            for (const name of possibleNames) {
                if (obj[name] && Array.isArray(obj[name])) {
                    return obj[name];
                }
            }
            
            // If JSON structure is nested (e.g. {data: {ExpenseCategories: [...]}})
            if (obj.data) {
                return this.getArrayFromObject(obj.data, possibleNames);
            }
            
            return []; // Default empty array
        },
        
        // Helper to ensure value is an array
        ensureArray(value) {
            if (Array.isArray(value)) return value;
            if (value === undefined || value === null) return [];
            // If value is a string like "[1, 2, 3]" (serialized array)
            if (typeof value === 'string' && value.trim().startsWith('[') && value.trim().endsWith(']')) {
                try {
                    return JSON.parse(value);
                } catch (e) {
                    console.warn('Failed to parse array from string:', value);
                    return [];
                }
            }
            return [value]; // Convert single value to array
        },
        
        fetchMasterDataFromAPI() {
            console.log('Fetching master data from API');
            
            fetch('/api/data-master')
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    return response.json();
                })
                .then(data => {
                    console.log('Data received from API:', data);
                    
                    // Try multiple property names
                    this.masterData = {
                        expenseCategories: this.getArrayFromObject(data, ['expenseCategories', 'ExpenseCategories']),
                        incomeCategories: this.getArrayFromObject(data, ['incomeCategories', 'IncomeCategories']),
                        paymentMethods: this.getArrayFromObject(data, ['paymentMethods', 'PaymentMethods']),
                        storageMedias: this.getArrayFromObject(data, ['storageMedias', 'StorageMedias', 'StorageMedia'])
                    };
                    
                    console.log('Master data loaded from API:', this.masterData);
                    
                    // Tampilkan notifikasi saat refresh selesai
                    if (this.loading) {
                        showToast('success', 'Data berhasil diperbarui');
                    }
                })
                .catch(error => {
                    console.error('Error fetching master data:', error);
                    showToast('error', 'Gagal memuat data master');
                })
                .finally(() => {
                    this.loading = false;
                });
        },
        
        // Helper method to get items based on type
        getItems(type) {
            switch (type) {
                case 'expense-categories':
                    return this.masterData.expenseCategories || [];
                case 'income-categories':
                    return this.masterData.incomeCategories || [];
                case 'payment-methods':
                    return this.masterData.paymentMethods || [];
                case 'storage-media':
                    return this.masterData.storageMedias || [];
                default:
                    return [];
            }
        }
    }));
});
