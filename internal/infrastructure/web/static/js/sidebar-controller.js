/**
 * Sidebar Controller - Menangani interaksi sidebar
 */
document.addEventListener('alpine:init', () => {
    Alpine.data('sidebarController', () => ({
        collapsed: false,
        mobileOpen: false,
        
        init() {
            // Check localStorage for saved state
            this.collapsed = localStorage.getItem('sidebarCollapsed') === 'true';
            
            // Add event listeners
            this.addEventListeners();
            
            // Handle resize events to reset mobile menu when screen size changes
            window.addEventListener('resize', () => {
                if (window.innerWidth > 768 && this.mobileOpen) {
                    this.mobileOpen = false;
                }
            });
        },
        
        addEventListeners() {
            // Listen for clicks outside when mobile menu is open
            document.addEventListener('click', (event) => {
                if (this.mobileOpen && !event.target.closest('.l-sidebar') && 
                    !event.target.closest('#mobileSidebarToggle')) {
                    this.mobileOpen = false;
                }
            });
        },
        
        toggleCollapse() {
            this.collapsed = !this.collapsed;
            localStorage.setItem('sidebarCollapsed', this.collapsed);
        },
        
        toggleMobileMenu() {
            this.mobileOpen = !this.mobileOpen;
        },
        
        closeMobileMenu() {
            this.mobileOpen = false;
        }
    }));
});
