/**
 * Contacts Application
 * Menangani fungsionalitas halaman kontak dan whitelist
 */
document.addEventListener('alpine:init', () => {
    Alpine.data('contactsApp', () => ({
        contacts: [],
        filteredContacts: [],
        searchQuery: '',
        filterType: 'all',
        useWhitelist: false,
        loading: true,
        
        // Modal states
        showAddModal: false,
        showEditModal: false,
        showDeleteModal: false,
        
        // Form models
        newContact: {
            name: '',
            phone: '',
            notes: '',
            whitelisted: false
        },
        editContact: {
            phone: '',
            name: '',
            notes: '',
            whitelisted: false
        },
        deleteContact: {
            phone: '',
            name: ''
        },
        
        initialize() {
            console.log('Initializing contacts app');
            
            // Load contacts
            this.fetchContacts();
            
            // Fetch whitelist settings
            this.fetchWhitelistSettings();
            
            // Watch for search queries and filter changes
            this.$watch('searchQuery', () => this.applyFilters());
            this.$watch('filterType', () => this.applyFilters());
        },
        
        fetchContacts() {
            this.loading = true;
            
            fetch('/api/contacts')
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    return response.json();
                })
                .then(data => {
                    if (Array.isArray(data.contacts)) {
                        this.contacts = data.contacts;
                        this.applyFilters();
                    } else {
                        console.error('Expected contacts array but got:', data);
                        this.contacts = [];
                        this.filteredContacts = [];
                    }
                    this.loading = false;
                })
                .catch(error => {
                    console.error('Error fetching contacts:', error);
                    this.showErrorToast('Gagal memuat daftar kontak');
                    this.contacts = [];
                    this.filteredContacts = [];
                    this.loading = false;
                });
        },
        
        fetchWhitelistSettings() {
            fetch('/api/whitelist/status')
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    return response.json();
                })
                .then(data => {
                    this.useWhitelist = data.enabled;
                    console.log('Whitelist status:', this.useWhitelist);
                })
                .catch(error => {
                    console.error('Error fetching whitelist settings:', error);
                });
        },
        
        applyFilters() {
            // Apply search filter
            let filtered = this.contacts;
            
            if (this.searchQuery) {
                const query = this.searchQuery.toLowerCase();
                filtered = filtered.filter(contact => 
                    contact.name.toLowerCase().includes(query) || 
                    contact.phone.toLowerCase().includes(query) ||
                    (contact.notes && contact.notes.toLowerCase().includes(query))
                );
            }
            
            // Apply type filter
            if (this.filterType === 'whitelist') {
                filtered = filtered.filter(contact => contact.whitelisted);
            } else if (this.filterType === 'regular') {
                filtered = filtered.filter(contact => !contact.whitelisted);
            }
            
            this.filteredContacts = filtered;
        },
        
        resetSearch() {
            this.searchQuery = '';
            this.filterType = 'all';
        },
        
        toggleWhitelist() {
            const newState = !this.useWhitelist;
            
            // Nonaktifkan tombol saat proses API
            const toggleButton = document.querySelector('.c-toggle-button');
            if (toggleButton) toggleButton.disabled = true;
            
            fetch('/api/whitelist/toggle', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    enabled: newState
                })
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    return response.json();
                })
                .then(data => {
                    if (data.success) {
                        // Update state setelah API berhasil
                        this.useWhitelist = data.enabled;
                        
                        // Gunakan jenis toast berbeda berdasarkan status whitelist
                        if (this.useWhitelist) {
                            // Aktif - gunakan toast success
                            this.showSuccessToast(`Mode whitelist diaktifkan`);
                        } else {
                            // Nonaktif - gunakan toast warning
                            this.showWarningToast(`Mode whitelist dinonaktifkan`);
                        }
                    } else {
                        throw new Error(data.error || 'Unknown error');
                    }
                })
                .catch(error => {
                    console.error('Error toggling whitelist:', error);
                    this.showErrorToast('Gagal mengubah pengaturan whitelist');
                })
                .finally(() => {
                    // Re-enable tombol
                    if (toggleButton) toggleButton.disabled = false;
                });
        },
        
        // Modal control functions
        openAddModal() {
            this.newContact = {
                name: '',
                phone: '',
                notes: '',
                whitelisted: false
            };
            this.showAddModal = true;
        },
        
        openEditModal(contact) {
            this.editContact = {
                phone: contact.phone,
                name: contact.name,
                notes: contact.notes || '',
                whitelisted: contact.whitelisted
            };
            this.showEditModal = true;
        },
        
        openDeleteModal(contact) {
            this.deleteContact = {
                phone: contact.phone,
                name: contact.name
            };
            this.showDeleteModal = true;
        },
        
        // CRUD operations
        addContact() {
            // Validate phone format
            if (!this.validatePhone(this.newContact.phone)) {
                this.showErrorToast('Format nomor telepon tidak valid');
                return;
            }
            
            // Validate required fields
            if (!this.newContact.name || !this.newContact.phone) {
                this.showErrorToast('Nama dan nomor telepon wajib diisi');
                return;
            }
            
            // Show loading state on button
            const submitBtn = document.activeElement;
            const originalText = submitBtn.innerHTML;
            submitBtn.disabled = true;
            submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Menyimpan...';
            
            fetch('/api/contacts/add', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    name: this.newContact.name,
                    phone: this.newContact.phone,
                    notes: this.newContact.notes,
                    isActive: this.newContact.whitelisted
                })
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    return response.json();
                })
                .then(data => {
                    if (data.success) {
                        this.showSuccessToast('Kontak berhasil ditambahkan');
                        this.showAddModal = false;
                        this.fetchContacts(); // Refresh the list
                    } else {
                        throw new Error(data.error || 'Unknown error');
                    }
                })
                .catch(error => {
                    console.error('Error adding contact:', error);
                    this.showErrorToast('Gagal menambahkan kontak');
                })
                .finally(() => {
                    // Reset button state
                    submitBtn.disabled = false;
                    submitBtn.innerHTML = originalText;
                });
        },
        
        updateContact() {
            // Validate required fields
            if (!this.editContact.name) {
                this.showErrorToast('Nama kontak wajib diisi');
                return;
            }
            
            // Show loading state
            const submitBtn = document.activeElement;
            const originalText = submitBtn.innerHTML;
            submitBtn.disabled = true;
            submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Menyimpan...';
            
            fetch('/api/contacts/update', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    phone: this.editContact.phone,
                    name: this.editContact.name,
                    notes: this.editContact.notes,
                    isActive: this.editContact.whitelisted
                })
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    return response.json();
                })
                .then(data => {
                    if (data.success) {
                        this.showSuccessToast('Kontak berhasil diperbarui');
                        this.showEditModal = false;
                        this.fetchContacts(); // Refresh the list
                    } else {
                        throw new Error(data.error || 'Unknown error');
                    }
                })
                .catch(error => {
                    console.error('Error updating contact:', error);
                    this.showErrorToast('Gagal memperbarui kontak');
                })
                .finally(() => {
                    // Reset button state
                    submitBtn.disabled = false;
                    submitBtn.innerHTML = originalText;
                });
        },
        
        deleteContact() {
            // Show loading state
            const deleteBtn = document.activeElement;
            const originalText = deleteBtn.innerHTML;
            deleteBtn.disabled = true;
            deleteBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Menghapus...';
            
            fetch('/api/contacts/delete', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    phone: this.deleteContact.phone
                })
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    return response.json();
                })
                .then(data => {
                    if (data.success) {
                        this.showSuccessToast('Kontak berhasil dihapus');
                        this.showDeleteModal = false;
                        this.fetchContacts(); // Refresh the list
                    } else {
                        throw new Error(data.error || 'Unknown error');
                    }
                })
                .catch(error => {
                    console.error('Error deleting contact:', error);
                    this.showErrorToast('Gagal menghapus kontak');
                })
                .finally(() => {
                    // Reset button state
                    deleteBtn.disabled = false;
                    deleteBtn.innerHTML = originalText;
                });
        },
        
        toggleWhitelistStatus(contact) {
            const newStatus = !contact.whitelisted;
            const contactEl = event.target.closest('.c-contact-card');
            
            if (contactEl) {
                contactEl.classList.add('c-contact-card--updating');
            }
            
            fetch('/api/whitelist/status', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    phone: contact.phone,
                    isActive: newStatus
                })
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    return response.json();
                })
                .then(data => {
                    if (data.success) {
                        // Update local state
                        const updatedContacts = this.contacts.map(c => {
                            if (c.phone === contact.phone) {
                                return { ...c, whitelisted: newStatus };
                            }
                            return c;
                        });
                        
                        this.contacts = updatedContacts;
                        this.applyFilters();
                        
                        const actionText = newStatus ? 'ditambahkan ke' : 'dihapus dari';
                        this.showSuccessToast(`Kontak ${actionText} whitelist`);
                    } else {
                        throw new Error(data.error || 'Unknown error');
                    }
                })
                .catch(error => {
                    console.error('Error updating whitelist status:', error);
                    this.showErrorToast('Gagal mengubah status whitelist');
                })
                .finally(() => {
                    if (contactEl) {
                        contactEl.classList.remove('c-contact-card--updating');
                    }
                });
        },
        
        // Helper functions
        validatePhone(phone) {
            // Allow numbers with optional + prefix
            return /^(\+)?[0-9]+$/.test(phone);
        },
        
        getInitials(name) {
            if (!name) return '?';
            
            const parts = name.trim().split(' ');
            if (parts.length === 1) {
                return parts[0].charAt(0).toUpperCase();
            }
            
            return (parts[0].charAt(0) + parts[parts.length - 1].charAt(0)).toUpperCase();
        },
        
        formatPhone(phone) {
            if (!phone) return '-';
            
            // Format as international number if not already formatted
            if (phone.startsWith('62')) {
                return '+' + phone;
            }
            
            return phone;
        },
        
        // Toast notifications
        showSuccessToast(message) {
            if (typeof showToast === 'function') {
                showToast('success', message);
            } else {
                alert(message);
            }
        },
        
        showErrorToast(message) {
            if (typeof showToast === 'function') {
                showToast('error', message);
            } else {
                alert('Error: ' + message);
            }
        },

        // Tambahkan metode toast warning
        showWarningToast(message) {
            if (typeof showToast === 'function') {
                showToast('warning', message);
            } else {
                alert(message);
            }
        }
    }));
});
