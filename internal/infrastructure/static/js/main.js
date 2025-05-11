/**
 * Main JavaScript file for Botopia dashboard
 */
document.addEventListener('DOMContentLoaded', function() {
    // Initialize mobile menu toggle functionality
    const mobileMenuBtn = document.getElementById('mobile-menu-button');
    const sidebar = document.getElementById('sidebar');
    const mainContent = document.querySelector('.main-content');

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
    
    // Close sidebar when clicking outside on mobile
    document.addEventListener('click', function(event) {
        if (
            sidebar && 
            window.innerWidth < 1024 && 
            !sidebar.contains(event.target) && 
            !mobileMenuBtn.contains(event.target) &&
            !sidebar.classList.contains('-translate-x-full')
        ) {
            closeSidebar();
        }
    });
    
    // Initialize current page in navigation
    highlightCurrentPage();
    
    // Add scroll event listener for subtle header effects
    const header = document.querySelector('.header');
    if (header && mainContent) {
        mainContent.addEventListener('scroll', function() {
            if (mainContent.scrollTop > 10) {
                header.classList.add('shadow-md');
                header.style.backgroundColor = 'rgba(15, 23, 42, 0.8)';
            } else {
                header.classList.remove('shadow-md');
                header.style.backgroundColor = '';
            }
        });
    }
});

function highlightCurrentPage() {
    // Get current page from URL path
    const path = window.location.pathname;
    
    // Find all nav links
    const navLinks = document.querySelectorAll('#nav-links a');
    
    navLinks.forEach(link => {
        // Remove active class from all links
        link.classList.remove('bg-primary-600', 'active');
        
        // Add active class to current page link
        if (
            (path === '/' && link.getAttribute('href') === '/dashboard') ||
            (link.getAttribute('href') === path)
        ) {
            link.classList.add('bg-primary-600', 'active');
        }
    });
}
