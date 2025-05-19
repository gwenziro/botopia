/**
 * Sidebar Controller
 * Menangani interaksi dan state sidebar
 */
document.addEventListener("alpine:init", () => {
  console.log("Registering sidebar component with Alpine...");

  // PENTING: Pastikan nama komponen 'sidebar' match dengan x-data="sidebar" di sidebar.html
  Alpine.data("sidebar", () => ({
    
    collapsed: false,
    mobileOpen: false,
    activePath: window.location.pathname,

    init() {
      console.log("Sidebar component initialized");

      // Cek local storage untuk state sidebar
      const storedState = localStorage.getItem("sidebar-collapsed");
      if (storedState === "true") {
        this.collapsed = true;
        document.body.classList.add("sidebar-collapsed");
        
        // Tambahkan class untuk main content
        const mainContent = document.querySelector('.l-main');
        if (mainContent) {
          mainContent.classList.add('l-side-bar-collapsed');
          console.log('Added l-side-bar-collapsed class to main content');
        }
        
        console.log("Setting collapsed state from localStorage: true");
      }

      // Handle resize event untuk responsive behavior
      window.addEventListener("resize", this.handleResize.bind(this));
      this.handleResize();

      // Tambahkan listener ke tombol collapse untuk debugging
      this.$nextTick(() => {
        // Periksa seluruh selector yang mungkin ada
        const collapseBtns = document.querySelector(".l-sidebar__collapse-btn");
        if (collapseBtns) {
          console.log("Sidebar collapse button found:", collapseBtns);
          console.log("Collapsed state:", this.collapsed);
          console.log("Body has sidebar-collapsed class:", document.body.classList.contains("sidebar-collapsed"));
        } else {
          console.warn("Sidebar collapse button not found in DOM");
        }
      });

      // Tambahkan deteksi path aktif
      this.activePath = window.location.pathname;
      
      // Delay untuk Alpine harusnya sudah dirender
      this.$nextTick(() => {
        this.highlightActivePath();
        
        // Dispatch resize event untuk memastikan layout diperbarui
        window.dispatchEvent(new Event('resize'));
      });
    },

    toggleCollapse() {
      console.log("toggleCollapse called");
      this.toggleSidebar();
    },

    toggleSidebar() {
      console.log("toggleSidebar called, current state:", this.collapsed);

      // Toggle state
      this.collapsed = !this.collapsed;

      // Update class pada body untuk styling sidebar
      if (this.collapsed) {
        document.body.classList.add("sidebar-collapsed");
        
        // Tambahkan class ke main content
        const mainContent = document.querySelector('.l-main');
        if (mainContent) {
          mainContent.classList.add('l-side-bar-collapsed');
          console.log('Added l-side-bar-collapsed class to main content');
        }
        
        console.log("Added sidebar-collapsed class to body");
      } else {
        document.body.classList.remove("sidebar-collapsed");
        
        // Hapus class dari main content
        const mainContent = document.querySelector('.l-main');
        if (mainContent) {
          mainContent.classList.remove('l-side-bar-collapsed');
          console.log('Removed l-side-bar-collapsed class from main content');
        }
        
        console.log("Removed sidebar-collapsed class from body");
      }

      // Simpan state ke local storage
      localStorage.setItem("sidebar-collapsed", this.collapsed);
      
      // Dispatch event untuk memastikan layout diperbarui
      window.dispatchEvent(new Event('resize'));
    },

    toggleMobileSidebar() {
      this.mobileOpen = !this.mobileOpen;

      // Update class pada body untuk mobile sidebar
      if (this.mobileOpen) {
        document.body.classList.add("sidebar-mobile-open");
      } else {
        document.body.classList.remove("sidebar-mobile-open");
      }
    },

    handleResize() {
      // Jika layar mobile dan sidebar terbuka, tutup otomatis
      if (window.innerWidth < 768 && this.mobileOpen) {
        this.mobileOpen = false;
        document.body.classList.remove("sidebar-mobile-open");
      }
    },

    closeMobileSidebar() {
      if (this.mobileOpen) {
        this.mobileOpen = false;
        document.body.classList.remove("sidebar-mobile-open");
      }
    },

    // Fungsi baru untuk memeriksa dan menerapkan highlight pada menu aktif
    highlightActivePath() {
      // Coba aktifkan juga melalui DOM jika Alpine binding tidak bekerja
      const menuItems = document.querySelectorAll('.l-sidebar__item');
      menuItems.forEach(item => {
        const href = item.getAttribute('href');
        if (href === this.activePath) {
          item.classList.add('active');
          console.log("Added active class to:", href);
        } else {
          item.classList.remove('active');
        }
      });
    },

    // Fungsi untuk memeriksa apakah path tertentu aktif
    isPathActive(path) {
      return this.activePath === path;
    },
  }));
});