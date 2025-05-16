/**
 * Main JavaScript file for Botopia dashboard
 */
document.addEventListener('DOMContentLoaded', function() {
    // Initialize mobile menu toggle functionality
    const mobileMenuBtn = document.getElementById('mobile-menu-button');
    const sidebar = document.getElementById('sidebar');
    
    if (mobileMenuBtn && sidebar) {
        mobileMenuBtn.addEventListener('click', function() {
            if (sidebar.classList.contains('-translate-x-full')) {
                sidebar.classList.remove('-translate-x-full');
                sidebar.classList.add('translate-x-0');
                
                // Add overlay on mobile
                const overlay = document.createElement('div');
                overlay.id = 'sidebar-overlay';
                overlay.className = 'fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden';
                overlay.addEventListener('click', closeSidebar);
                document.body.appendChild(overlay);
            } else {
                closeSidebar();
            }
        });
    }
    
    // Function to close sidebar
    function closeSidebar() {
        if (sidebar) {
            sidebar.classList.remove('translate-x-0');
            sidebar.classList.add('-translate-x-full');
            
            // Remove overlay
            const overlay = document.getElementById('sidebar-overlay');
            if (overlay) {
                overlay.remove();
            }
        }
    }
    
    // Initialize current page in navigation
    highlightCurrentPage();
    
    // Add scroll event listener for header effects
    const header = document.querySelector('.header');
    const mainNav = document.getElementById('main-nav');
    
    if (header) {
        window.addEventListener('scroll', function() {
            if (window.scrollY > 10) {
                header.classList.add('shadow-md');
                header.style.backgroundColor = 'rgba(15, 23, 42, 0.8)';
            } else {
                header.classList.remove('shadow-md');
                header.style.backgroundColor = '';
            }
        });
    }
    
    // Sticky navigation for homepage
    if (mainNav) {
        window.addEventListener('scroll', function() {
            if (window.scrollY > 50) {
                mainNav.classList.add('scrolled');
            } else {
                mainNav.classList.remove('scrolled');
            }
        });
    }
});

function highlightCurrentPage() {
    // Get current page path
    const path = window.location.pathname;
    
    // Find all nav links
    const navLinks = document.querySelectorAll('#sidebar .nav-link');
    
    // Remove active class from all links
    navLinks.forEach(link => {
        link.classList.remove('active');
        
        // Add active class to current page link
        const href = link.getAttribute('href');
        if (href === path || 
            (path === '/' && href === '/dashboard') ||
            (href !== '/' && path.startsWith(href))) {
            link.classList.add('active');
        }
    });
}
