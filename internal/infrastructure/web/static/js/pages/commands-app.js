/**
 * Commands Application
 * 
 * Menampilkan dan mengelola daftar commands untuk bot WhatsApp
 */
document.addEventListener('alpine:init', () => {
    Alpine.data('commandsApp', () => ({
        loading: true,
        commands: [],
        filteredCommands: [],
        categories: [],
        totalCommands: 0,
        searchQuery: '',
        categoryFilter: '',
        debug: false,
        
        initCommands() {
            console.log('Initializing commands app from alpine:init');
            
            // Coba ambil parameter filter dari URL
            const urlParams = new URLSearchParams(window.location.search);
            if (urlParams.has('q')) {
                this.searchQuery = urlParams.get('q');
            }
            if (urlParams.has('category')) {
                this.categoryFilter = urlParams.get('category');
            }
            
            // Pantau perubahan filter untuk update URL
            this.$watch('searchQuery', () => {
                this.filterCommands();
                this.updateUrlParams();
            });
            
            this.$watch('categoryFilter', () => {
                this.filterCommands();
                this.updateUrlParams();
            });
            
            // Load data commands
            this.loadCommands();
        },
        
        loadCommands() {
            this.loading = true;
            
            // Strategi 1: Ambil dari script tag JSON
            const commandsDataElement = document.getElementById('commands-data');
            if (commandsDataElement && commandsDataElement.textContent.trim()) {
                try {
                    console.log('Loading commands from embedded JSON data');
                    const jsonText = commandsDataElement.textContent.trim();
                    console.log('JSON text length:', jsonText.length);
                    
                    // PERBAIKAN: Handle double-encoded JSON - decode sekali untuk dapatkan inner JSON string
                    let jsonData;
                    try {
                        const outerParsed = JSON.parse(jsonText);
                        console.log('First parse result (likely still a string):', typeof outerParsed);
                        
                        // Cek apakah hasil parsing pertama adalah string (double-encoded)
                        if (typeof outerParsed === 'string') {
                            // Parse lagi untuk mendapatkan object sebenarnya
                            jsonData = JSON.parse(outerParsed);
                            console.log('Second parse needed - final result:', typeof jsonData);
                        } else {
                            // Sudah dalam bentuk object
                            jsonData = outerParsed;
                        }
                    } catch (innerError) {
                        console.error('Error during JSON double-parsing:', innerError);
                        jsonData = {};
                    }
                    
                    console.log('Final parsed commands data:', jsonData);
                    this.processCommandsData(jsonData);
                    return;
                } catch (error) {
                    console.error('Error parsing JSON from script tag:', error);
                }
            }
            
            // Strategi 2: Ambil dari API
            console.log('Fetching commands from API');
            fetch('/api/commands')
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    return response.json();
                })
                .then(data => {
                    console.log('Received commands data from API:', data);
                    this.processCommandsData(data);
                })
                .catch(error => {
                    console.error('Error fetching commands:', error);
                    this.commands = [];
                    this.filteredCommands = [];
                    this.categories = [];
                    this.totalCommands = 0;
                    this.loading = false;
                    
                    if (typeof showToast === 'function') {
                        showToast('error', 'Gagal memuat data commands');
                    }
                });
        },
        
        processCommandsData(data) {
            console.log('Processing commands data');
            
            // PERBAIKAN: Verifikasi bahwa data adalah object sebelum diproses
            if (!data || typeof data !== 'object') {
                console.error('Invalid data format, expected object but got:', typeof data);
                this.loading = false;
                return;
            }
            
            // Reset array commands terlebih dahulu
            this.commands = [];
            
            // Konversi dari object ke array dengan cara yang benar
            Object.keys(data).forEach(key => {
                const cmd = data[key];
                if (cmd && typeof cmd === 'object') {
                    this.commands.push({
                        name: key,
                        description: cmd.description || 'Tidak ada deskripsi',
                        category: cmd.category || 'Umum',
                        usage: cmd.usage || `!${key}`
                    });
                }
            });
            
            console.log('Commands array created:', this.commands);
            console.log('Number of commands:', this.commands.length);
            
            // Sort by name
            this.commands.sort((a, b) => a.name.localeCompare(b.name));
            
            // Extract categories
            const categorySet = new Set();
            this.commands.forEach(cmd => {
                if (cmd.category) {
                    categorySet.add(cmd.category);
                }
            });
            this.categories = Array.from(categorySet).sort();
            
            this.totalCommands = this.commands.length;
            console.log(`Processed ${this.totalCommands} commands with ${this.categories.length} categories`);
            
            // Apply initial filters if any
            this.filterCommands();
            
            this.loading = false;
        },
        
        filterCommands() {
            console.log(`Filtering commands. Query: "${this.searchQuery}", Category: "${this.categoryFilter}"`);
            
            // Filter berdasarkan search query dan kategori
            this.filteredCommands = this.commands.filter(cmd => {
                const matchesSearch = !this.searchQuery || 
                    cmd.name.toLowerCase().includes(this.searchQuery.toLowerCase()) || 
                    cmd.description.toLowerCase().includes(this.searchQuery.toLowerCase());
                    
                const matchesCategory = !this.categoryFilter || 
                    cmd.category === this.categoryFilter;
                    
                return matchesSearch && matchesCategory;
            });
            
            console.log(`Filtered to ${this.filteredCommands.length} commands`);
        },
        
        resetFilters() {
            console.log('Resetting filters');
            this.searchQuery = '';
            this.categoryFilter = '';
            this.filterCommands();
            this.updateUrlParams();
        },
        
        updateUrlParams() {
            const url = new URL(window.location);
            
            if (this.searchQuery) {
                url.searchParams.set('q', this.searchQuery);
            } else {
                url.searchParams.delete('q');
            }
            
            if (this.categoryFilter) {
                url.searchParams.set('category', this.categoryFilter);
            } else {
                url.searchParams.delete('category');
            }
            
            history.replaceState({}, '', url);
        }
    }));
});

// Tambahan fallback untuk memastikan Alpine.js memuat komponen
document.addEventListener('DOMContentLoaded', function() {
    console.log('Commands app: DOM fully loaded');
    
    // Cek apakah Alpine.js sudah tersedia
    if (typeof Alpine !== 'undefined') {
        console.log('Alpine commandsApp component registered successfully');
    } else {
        console.warn('Alpine.js is not available yet. Will try to register component when it becomes available.');
        
        // Coba lagi setelah jeda pendek
        setTimeout(() => {
            if (typeof Alpine !== 'undefined') {
                console.log('Alpine.js is now available, registering command component');
                
                // Ulangi pendaftaran komponen jika perlu
                if (!Alpine.__component_defs || !Alpine.__component_defs.commandsApp) {
                    Alpine.data('commandsApp', function() {
                        // Implementasi sesuai dengan yang di atas
                        // ...
                    });
                }
            }
        }, 500);
    }
});
