/**
 * Contacts Application
 * Menangani fungsionalitas halaman kontak
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
            fetch('/api/contacts')
                .then(response => response.json())
                .then(data => {
                    if (Array.isArray(data.contacts)) {
                        this.contacts = data.contacts;
                        this.filteredContacts = [...this.contacts];
                    }
                    this.loading = false;
                })
                .catch(error => {
                    console.error('Error fetching contacts:', error);
                    showToast('error', 'Gagal memuat daftar kontak');
                    this.loading = false;
                });
                
            // Fetch whitelisted contacts separately
            fetch('/api/contacts/whitelist')
                .then(response => response.json())
                .then(data => {
                    if (Array.isArray(data.contacts)) {
                        this.whitelistedContacts = data.contacts;
                    }
                })
                .catch(error => {
                    console.error('Error fetching whitelisted contacts:', error);
                });
        },
        
        fetchWhitelistSettings() {
            fetch('/api/whitelist/status')
                .then(response => response.json())
                .then(data => {
                    this.useWhitelist = data.enabled;
                })
                .catch(error => {
                    console.error('Error fetching whitelist settings:', error);
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
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        this.useWhitelist = data.enabled;
                        showToast('success', `Whitelist ${this.useWhitelist ? 'diaktifkan' : 'dinonaktifkan'}`);
                    }
                })
                .catch(error => {
                    console.error('Error toggling whitelist:', error);
                    showToast('error', 'Gagal mengubah pengaturan whitelist');
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
                whitelisted: false
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
                whitelisted: contact.whitelisted
            };
            this.showEditModal = true;
        },
        
        openDeleteModal(contact) {
            this.currentContact = contact;
            this.showDeleteModal = true;
        },
        
        addContact() {
            // Validate phone format
            if (!this.validatePhone(this.newContact.phone)) {
                showToast('error', 'Format nomor telepon tidak valid');
                return;
            }
            
            fetch('/api/contacts/add', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(this.newContact)
            })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        showToast('success', 'Kontak berhasil ditambahkan');
                        this.showAddModal = false;
                        this.fetchContacts();
                    } else {
                        throw new Error(data.error || 'Failed to add contact');
                    }
                })
                .catch(error => {
                    console.error('Error adding contact:', error);
                    showToast('error', 'Gagal menambahkan kontak');
                });
        },
        
        updateContact() {
            if (!this.currentContact || !this.currentContact.id) {
                showToast('error', 'Data kontak tidak valid');
                return;
            }
            
            // Validate phone format
            if (!this.validatePhone(this.newContact.phone)) {
                showToast('error', 'Format nomor telepon tidak valid');
                return;
            }
            
            fetch('/api/contacts/update', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    id: this.currentContact.id,
                    ...this.newContact
                })
            })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        showToast('success', 'Kontak berhasil diperbarui');
                        this.showEditModal = false;
                        this.fetchContacts();
                    } else {
                        throw new Error(data.error || 'Failed to update contact');
                    }
                })
                .catch(error => {
                    console.error('Error updating contact:', error);
                    showToast('error', 'Gagal memperbarui kontak');
                });
        },
        
        deleteContact() {
            if (!this.currentContact || !this.currentContact.id) {
                showToast('error', 'Data kontak tidak valid');
                return;
            }
            
            fetch('/api/contacts/delete', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    id: this.currentContact.id
                })
            })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        showToast('success', 'Kontak berhasil dihapus');
                        this.showDeleteModal = false;
                        this.fetchContacts();
                    } else {
                        throw new Error(data.error || 'Failed to delete contact');
                    }
                })
                .catch(error => {
                    console.error('Error deleting contact:', error);
                    showToast('error', 'Gagal menghapus kontak');
                });
        },
        
        toggleWhitelistStatus(contact) {
            fetch('/api/whitelist/status', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    id: contact.id,
                    whitelisted: !contact.whitelisted
                })
            })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        // Update local data
                        this.contacts = this.contacts.map(c => {
                            if (c.id === contact.id) {
                                return { ...c, whitelisted: !c.whitelisted };
                            }
                            return c;
                        });
                        this.filterContacts();
                        showToast('success', `Kontak ${!contact.whitelisted ? 'ditambahkan ke' : 'dihapus dari'} whitelist`);
                    } else {
                        throw new Error(data.error || 'Failed to update whitelist status');
                    }
                })
                .catch(error => {
                    console.error('Error updating whitelist status:', error);
                    showToast('error', 'Gagal mengubah status whitelist');
                });
        },
        
        validatePhone(phone) {
            // Validate phone number format (number with optional + prefix)
            return /^(\+)?[\d]+$/.test(phone);
        }
    }));
});
