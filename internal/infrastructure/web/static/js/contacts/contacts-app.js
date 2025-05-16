/**
 * Contacts App
 * Aplikasi untuk mengelola kontak dan whitelist
 */
document.addEventListener('alpine:init', () => {
    Alpine.data('contactsApp', () => ({
        activeTab: 'all',
        loading: false,
        contacts: [],
        whitelistedContacts: [],
        useWhitelist: false,
        showContactModal: false,
        editMode: false,
        isSaving: false,
        contactForm: {
            phone: '',
            name: '',
            notes: '',
            isActive: true
        },
        
        initContacts() {
            console.log('Initializing contacts app');
            this.loading = true;
            
            // Ambil data kontak dan status whitelist
            Promise.all([
                this.fetchContacts(),
                this.fetchWhitelistStatus()
            ]).finally(() => {
                this.loading = false;
            });
            
            // Set active tab from URL if available
            const urlParams = new URLSearchParams(window.location.search);
            if (urlParams.has('tab')) {
                this.activeTab = urlParams.get('tab');
            }
            
            // Watch for tab changes to update URL
            this.$watch('activeTab', (value) => {
                const url = new URL(window.location);
                url.searchParams.set('tab', value);
                history.replaceState(null, '', url);
            });
        },
        
        refreshContacts() {
            this.loading = true;
            
            Promise.all([
                this.fetchContacts(),
                this.fetchWhitelistStatus()
            ]).then(() => {
                showToast('success', 'Data kontak berhasil diperbarui');
            }).catch(error => {
                showToast('error', 'Gagal memperbarui data: ' + error.message);
            }).finally(() => {
                this.loading = false;
            });
        },
        
        fetchContacts() {
            return fetch('/api/contacts')
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Failed to fetch contacts');
                    }
                    return response.json();
                })
                .then(data => {
                    this.contacts = data.contacts || [];
                    // Filter untuk whitelist
                    this.updateWhitelistedContacts();
                })
                .catch(error => {
                    console.error('Error fetching contacts:', error);
                    showToast('error', 'Gagal memuat kontak');
                });
        },
        
        fetchWhitelistStatus() {
            return fetch('/api/whitelist/status')
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Failed to fetch whitelist status');
                    }
                    return response.json();
                })
                .then(data => {
                    this.useWhitelist = data.enabled;
                })
                .catch(error => {
                    console.error('Error fetching whitelist status:', error);
                });
        },
        
        updateWhitelistedContacts() {
            // Filter kontak yang aktif (dalam whitelist)
            this.whitelistedContacts = this.contacts.filter(contact => contact.isActive);
        },
        
        toggleWhitelist() {
            fetch('/api/whitelist/toggle', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ enabled: this.useWhitelist })
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Failed to update whitelist status');
                    }
                    return response.json();
                })
                .then(data => {
                    showToast('success', this.useWhitelist ? 
                        'Whitelist diaktifkan - Bot hanya merespon kontak yang terdaftar' : 
                        'Whitelist dinonaktifkan - Bot merespon semua kontak');
                })
                .catch(error => {
                    console.error('Error toggling whitelist:', error);
                    showToast('error', 'Gagal mengubah status whitelist');
                    this.useWhitelist = !this.useWhitelist; // Revert change on error
                });
        },
        
        openAddModal() {
            console.log('Opening add modal');
            this.editMode = false;
            this.contactForm = {
                phone: '',
                name: '',
                notes: '',
                isActive: true
            };
            this.showContactModal = true;
        },
        
        editContact(contact) {
            this.editMode = true;
            this.contactForm = {
                phone: contact.phone,
                name: contact.name,
                notes: contact.notes,
                isActive: contact.isActive
            };
            this.showContactModal = true;
        },
        
        saveContact() {
            // Validasi nomor telepon
            if (!this.contactForm.phone) {
                showToast('error', 'Nomor telepon tidak boleh kosong');
                return;
            }
            
            this.isSaving = true;
            
            const endpoint = this.editMode ? '/api/contacts/update' : '/api/contacts/add';
            
            fetch(endpoint, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(this.contactForm)
            })
                .then(response => {
                    if (!response.ok) {
                        return response.json().then(data => {
                            throw new Error(data.error || 'Failed to save contact');
                        });
                    }
                    return response.json();
                })
                .then(data => {
                    showToast('success', this.editMode ? 
                        'Kontak berhasil diperbarui' : 
                        'Kontak baru berhasil ditambahkan');
                    this.showContactModal = false;
                    this.refreshContacts();
                })
                .catch(error => {
                    console.error('Error saving contact:', error);
                    showToast('error', error.message || 'Gagal menyimpan kontak');
                })
                .finally(() => {
                    this.isSaving = false;
                });
        },
        
        deleteContact(contact) {
            if (!confirm(`Yakin ingin menghapus kontak ${contact.name || contact.phone}?`)) {
                return;
            }
            
            fetch('/api/contacts/delete', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ phone: contact.phone })
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Failed to delete contact');
                    }
                    return response.json();
                })
                .then(data => {
                    showToast('success', 'Kontak berhasil dihapus');
                    this.refreshContacts();
                })
                .catch(error => {
                    console.error('Error deleting contact:', error);
                    showToast('error', 'Gagal menghapus kontak');
                });
        },
        
        toggleContactStatus(contact) {
            const newStatus = !contact.isActive;
            const actionName = newStatus ? 'menambahkan ke whitelist' : 'menghapus dari whitelist';
            
            fetch('/api/whitelist/status', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    phone: contact.phone,
                    isActive: newStatus
                })
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error(`Gagal ${actionName}`);
                    }
                    return response.json();
                })
                .then(data => {
                    // Update local contact status
                    contact.isActive = newStatus;
                    this.updateWhitelistedContacts();
                    
                    showToast('success', newStatus ? 
                        `${contact.name || contact.phone} ditambahkan ke whitelist` : 
                        `${contact.name || contact.phone} dihapus dari whitelist`);
                })
                .catch(error => {
                    console.error('Error toggling contact status:', error);
                    showToast('error', error.message || 'Gagal mengubah status kontak');
                });
        },
        
        removeFromWhitelist(contact) {
            this.toggleContactStatus(contact);
        }
    }));
});
