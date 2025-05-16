/**
 * Sidebar menu activation functionality
 * 
 * Handles setting the active class on the current page's menu item
 */
document.addEventListener('DOMContentLoaded', function() {
    // Get current path
    const currentPath = window.location.pathname;
    
    // Find all nav links in sidebar
    const navLinks = document.querySelectorAll('.sidebar .nav-link');
    
    // Remove any existing active class
    navLinks.forEach(link => {
        link.classList.remove('active');
    });
    
    // Set active class on matching link
    navLinks.forEach(link => {
        const href = link.getAttribute('href');
        
        // For exact matches (like /dashboard)
        if (href === currentPath) {
            link.classList.add('active');
        }
        
        // For partial matches (like /data-master/categories would match /data-master)
        else if (href !== '/' && currentPath.startsWith(href)) {
            link.classList.add('active');
        }
    });
    
    // Special case for dashboard being home
    if (currentPath === '/' || currentPath === '') {
        const dashboardLink = document.querySelector('.sidebar .nav-link[href="/dashboard"]');
        if (dashboardLink) {
            dashboardLink.classList.add('active');
        }
    }
});
