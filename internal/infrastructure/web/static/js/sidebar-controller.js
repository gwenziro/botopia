/**
 * Sidebar Controller
 * 
 * Mengelola keadaan dan perilaku sidebar
 */
document.addEventListener('alpine:init', () => {
    Alpine.data('sidebarController', () => ({
        sidebarCollapsed: false,
        
        init() {
            // Cek local storage untuk status sidebar sebelumnya
            const savedState = localStorage.getItem('sidebarCollapsed');
            if (savedState !== null) {
                this.sidebarCollapsed = savedState === 'true';
            }
            
            // Tambahkan data-title ke semua nav-link untuk tooltip
            this.addTitlesToNavLinks();
            
            // Tambahkan listener untuk perubahan ukuran layar
            this.setupResizeListener();
        },
        
        toggleCollapse() {
            this.sidebarCollapsed = !this.sidebarCollapsed;
            localStorage.setItem('sidebarCollapsed', this.sidebarCollapsed);
        },
        
        addTitlesToNavLinks() {
            // Tambahkan atribut data-title ke semua nav-link berdasarkan teks
            document.querySelectorAll('.nav-link').forEach(link => {
                const textElement = link.querySelector('.sidebar-text');
                if (textElement) {
                    const text = textElement.innerText;
                    link.setAttribute('data-title', text);
                }
            });
        },
        
        setupResizeListener() {
            // Otomatis collapse sidebar pada layar kecil
            const mobileBreakpoint = 768; // Batas ukuran mobile dalam pixel
            
            const handleResize = () => {
                if (window.innerWidth < mobileBreakpoint) {
                    this.sidebarCollapsed = true;
                }
            };
            
            // Jalankan sekali saat init
            handleResize();
            
            // Tambahkan listener untuk resize
            window.addEventListener('resize', handleResize);
        }
    }));
});
