/**
 * Contacts Application
 * Menangani fungsionalitas halaman kontak dan whitelist
 */
document.addEventListener('alpine:init', () => {
    Alpine.data('contactApp', () => ({
        contacts: [],
        whitelistedContacts: [],
        searchQuery: '',
        filteredContacts: [],
        useWhitelist: false,
        loading: true,
        currentContact: null,
        showAddModal: false,
        showEditModal: false,
        showDeleteModal: false,
        newContact: {
            name: '',
            phone: '',
            whitelisted: false
        },
        
        initialize() {
            console.log('Initializing contact app');
            
            // Load contacts
            this.fetchContacts();
            
            // Fetch whitelist settings
            this.fetchWhitelistSettings();
            
            // Watch for search queries
            this.$watch('searchQuery', () => this.filterContacts());
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
                        this.filteredContacts = [...this.contacts];
                    } else {
                        // Fallback jika struktur response tidak sesuai ekspektasi
                        this.contacts = [];
                        this.filteredContacts = [];
                        console.warn('Unexpected API response format', data);
                    }
                    this.loading = false;
                })
                .catch(error => {
                    console.error('Error fetching contacts:', error);
                    if (typeof showToast === 'function') {
                        showToast('error', 'Gagal memuat daftar kontak');
                    }
                    this.contacts = [];
                    this.filteredContacts = [];
                    this.loading = false;
                });
                
            // Fetch whitelisted contacts separately
            fetch('/api/contacts/whitelist')
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    return response.json();
                })
                .then(data => {
                    if (Array.isArray(data.contacts)) {
                        this.whitelistedContacts = data.contacts;
                    } else {
                        this.whitelistedContacts = [];
                        console.warn('Unexpected API response format for whitelist', data);
                    }
                })
                .catch(error => {
                    console.error('Error fetching whitelisted contacts:', error);
                    this.whitelistedContacts = [];
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
                })
                .catch(error => {
                    console.error('Error fetching whitelist settings:', error);
                    // Default ke false jika gagal
                    this.useWhitelist = false;
                });
        },
        
        toggleWhitelist() {
            fetch('/api/whitelist/toggle', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    enabled: !this.useWhitelist
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
                        this.useWhitelist = data.enabled;
                        
                        if (typeof showToast === 'function') {
                            showToast(
                                'success', 
                                `Mode whitelist ${this.useWhitelist ? 'diaktifkan' : 'dinonaktifkan'}`
                            );
                        }
                    } else {
                        throw new Error(data.error || 'Failed to toggle whitelist');
                    }
                })
                .catch(error => {
                    console.error('Error toggling whitelist:', error);
                    
                    if (typeof showToast === 'function') {
                        showToast('error', 'Gagal mengubah pengaturan whitelist');
                    }
                });
        },
        
        filterContacts() {
            if (!this.searchQuery.trim()) {
                this.filteredContacts = [...this.contacts];
                return;
            }
            
            const query = this.searchQuery.toLowerCase();
            this.filteredContacts = this.contacts.filter(contact =>
                contact.name.toLowerCase().includes(query) || 
                contact.phone.toLowerCase().includes(query)
            );
        },
        
        resetForm() {
            this.newContact = {
                name: '',
                phone: '',
                whitelisted: false,
                notes: ''
            };
        },
        
        openAddModal() {
            this.resetForm();
            this.showAddModal = true;
        },
        
        openEditModal(contact) {
            this.currentContact = contact;
            this.newContact = {
                name: contact.name,
                phone: contact.phone,
                whitelisted: contact.whitelisted,
                notes: contact.notes || ''
            };
            this.showEditModal = true;
        },
        
        openDeleteModal(contact) {
            this.currentContact = contact;
            this.showDeleteModal = true;
        },
        
        addContact() {
            // Validasi format nomor telepon
            if (!this.validatePhone(this.newContact.phone)) {
                if (typeof showToast === 'function') {
                    showToast('error', 'Format nomor telepon tidak valid');
                }
                return;
            }
            
            // Normalisasi format telepon
            const phone = this.normalizePhone(this.newContact.phone);
            
            fetch('/api/contacts/add', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    name: this.newContact.name,
                    phone: phone,
                    notes: this.newContact.notes || '',
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
                        if (typeof showToast === 'function') {
                            showToast('success', 'Kontak berhasil ditambahkan');
                        }
                        
                        this.showAddModal = false;
                        this.fetchContacts(); // Reload contacts
                    } else {
                        throw new Error(data.error || 'Failed to add contact');
                    }
                })
                .catch(error => {
                    console.error('Error adding contact:', error);
                    
                    if (typeof showToast === 'function') {
                        showToast('error', 'Gagal menambahkan kontak: ' + error.message);
                    }
                });
        },
        
        updateContact() {
            if (!this.currentContact) {
                if (typeof showToast === 'function') {
                    showToast('error', 'Data kontak tidak valid');
                }
                return;
            }
            
            fetch('/api/contacts/update', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    phone: this.currentContact.phone,
                    name: this.newContact.name,
                    notes: this.newContact.notes || '',
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
                        if (typeof showToast === 'function') {
                            showToast('success', 'Kontak berhasil diperbarui');
                        }
                        
                        this.showEditModal = false;
                        this.fetchContacts(); // Reload contacts
                    } else {
                        throw new Error(data.error || 'Failed to update contact');
                    }
                })
                .catch(error => {
                    console.error('Error updating contact:', error);
                    
                    if (typeof showToast === 'function') {
                        showToast('error', 'Gagal memperbarui kontak: ' + error.message);
                    }
                });
        },
        
        deleteContact() {
            if (!this.currentContact) {
                if (typeof showToast === 'function') {
                    showToast('error', 'Data kontak tidak valid');
                }
                return;
            }
            
            fetch('/api/contacts/delete', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    phone: this.currentContact.phone
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
                        if (typeof showToast === 'function') {
                            showToast('success', 'Kontak berhasil dihapus');
                        }
                        
                        this.showDeleteModal = false;
                        this.fetchContacts(); // Reload contacts
                    } else {
                        throw new Error(data.error || 'Failed to delete contact');
                    }
                })
                .catch(error => {
                    console.error('Error deleting contact:', error);
                    
                    if (typeof showToast === 'function') {
                        showToast('error', 'Gagal menghapus kontak: ' + error.message);
                    }
                });
        },
        
        toggleWhitelistStatus(contact) {
            fetch('/api/whitelist/status', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    phone: contact.phone,
                    isActive: !contact.whitelisted
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
                        // Update contact in local array
                        this.contacts = this.contacts.map(c => {
                            if (c.phone === contact.phone) {
                                return { ...c, whitelisted: !c.whitelisted };
                            }
                            return c;
                        });
                        
                        // Update filtered contacts too
                        this.filteredContacts = this.filteredContacts.map(c => {
                            if (c.phone === contact.phone) {
                                return { ...c, whitelisted: !c.whitelisted };
                            }
                            return c;
                        });
                        
                        // Update whitelisted contacts list
                        this.fetchContacts();
                        
                        if (typeof showToast === 'function') {
                            const actionText = !contact.whitelisted ? 
                                'ditambahkan ke' : 'dihapus dari';
                            showToast('success', `Kontak ${actionText} whitelist`);
                        }
                    } else {
                        throw new Error(data.error || 'Failed to update whitelist status');
                    }
                })
                .catch(error => {
                    console.error('Error updating whitelist status:', error);
                    
                    if (typeof showToast === 'function') {
                        showToast('error', 'Gagal mengubah status whitelist');
                    }
                });
        },
        
        validatePhone(phone) {
            // Validasi format nomor telepon, menerima +62, 62, atau 08
            return /^(\+?62|0)[0-9]{9,13}$/.test(phone);
        },
        
        normalizePhone(phone) {
            // Normalisasi format nomor telepon
            if (phone.startsWith('0')) {
                return '+62' + phone.substring(1);
            }
            
            if (phone.startsWith('62')) {
                return '+' + phone;
            }
            
            return phone;
        }
    }));
});
