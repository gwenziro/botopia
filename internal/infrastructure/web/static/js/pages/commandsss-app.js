/**
 * Commands Application
 * Menangani fungsionalitas halaman commands
 */
document.addEventListener('alpine:init', () => {
    Alpine.data('commandsApp', () => ({
        commands: [],
        filteredCommands: [],
        categories: [],
        categoryFilter: '',
        searchQuery: '',
        viewMode: 'grid',
        loading: true,
        
        initializeCommands() {
            console.log('Initializing commands app');
            
            // Load data
            this.fetchCommands();
            
            // Watch for search and category filter changes
            this.$watch('searchQuery', () => this.filterCommands());
            this.$watch('categoryFilter', () => this.filterCommands());
        },
        
        fetchCommands() {
            // Check if commands are available in global object first
            if (window.botopiaCommands) {
                this.processCommands(window.botopiaCommands);
                return;
            }
            
            // Fetch from API if not available
            fetch('/api/commands')
                .then(response => response.json())
                .then(data => {
                    this.processCommands(data);
                })
                .catch(error => {
                    console.error('Error loading commands:', error);
                    showToast('error', 'Gagal memuat daftar command');
                    this.loading = false;
                });
        },
        
        processCommands(data) {
            // Convert commands object to array
            this.commands = Object.keys(data).map(key => {
                return {
                    name: key,
                    description: data[key].description || 'Tidak ada deskripsi',
                    category: data[key].category || 'Umum',
                    usage: data[key].usage || `!${key}`
                };
            });
            
            // Sort commands alphabetically
            this.commands.sort((a, b) => a.name.localeCompare(b.name));
            
            // Extract unique categories
            this.categories = [...new Set(this.commands.map(cmd => cmd.category))].sort();
            
            // Initial filtering
            this.filterCommands();
            
            this.loading = false;
        },
        
        filterCommands() {
            let result = [...this.commands];
            
            // Apply category filter
            if (this.categoryFilter) {
                result = result.filter(cmd => cmd.category === this.categoryFilter);
            }
            
            // Apply search query
            if (this.searchQuery.trim() !== '') {
                const query = this.searchQuery.toLowerCase();
                result = result.filter(cmd => 
                    cmd.name.toLowerCase().includes(query) || 
                    cmd.description.toLowerCase().includes(query)
                );
            }
            
            this.filteredCommands = result;
        },
        
        getCommandIcon(name) {
            // Mapping ikon untuk command berdasarkan nama atau kategori
            const iconMap = {
                'panduan': 'fa-info-circle',
                'ping': 'fa-heart-pulse',
                'keluar': 'fa-arrow-right-from-bracket',
                'masuk': 'fa-arrow-right-to-bracket',
                'unggah': 'fa-upload'
            };
            
            return iconMap[name] || 'fa-terminal';
        }
    }));
});
