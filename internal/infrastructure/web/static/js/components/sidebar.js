/**
 * Sidebar Controller
 * Menangani interaksi dan state sidebar
 */
document.addEventListener('alpine:init', () => {
    Alpine.data('sidebar', () => ({
        isCollapsed: false,
        isMobileOpen: false,
        
        init() {
            // Cek local storage untuk state sidebar
            const storedState = localStorage.getItem('sidebar-collapsed');
            if (storedState === 'true') {
                this.isCollapsed = true;
                document.body.classList.add('sidebar-collapsed');
            }
            
            // Handle resize event untuk responsive behavior
            window.addEventListener('resize', this.handleResize.bind(this));
            this.handleResize();
        },
        
        toggleSidebar() {
            this.isCollapsed = !this.isCollapsed;
            
            // Update class pada body untuk styling sidebar
            if (this.isCollapsed) {
                document.body.classList.add('sidebar-collapsed');
            } else {
                document.body.classList.remove('sidebar-collapsed');
            }
            
            // Simpan state ke local storage
            localStorage.setItem('sidebar-collapsed', this.isCollapsed);
            
            // Dispatch event untuk komponen lain
            window.dispatchEvent(new CustomEvent('sidebar-toggle', {
                detail: { isCollapsed: this.isCollapsed }
            }));
        },
        
        toggleMobileSidebar() {
            this.isMobileOpen = !this.isMobileOpen;
            
            // Update class pada body untuk mobile sidebar
            if (this.isMobileOpen) {
                document.body.classList.add('sidebar-mobile-open');
            } else {
                document.body.classList.remove('sidebar-mobile-open');
            }
        },
        
        handleResize() {
            // Jika layar mobile dan sidebar terbuka, tutup otomatis
            if (window.innerWidth < 768 && this.isMobileOpen) {
                this.isMobileOpen = false;
                document.body.classList.remove('sidebar-mobile-open');
            }
        },
        
        closeMobileSidebar() {
            if (this.isMobileOpen) {
                this.isMobileOpen = false;
                document.body.classList.remove('sidebar-mobile-open');
            }
        }
    }));
});
