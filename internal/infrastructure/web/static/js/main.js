/**
 * Main JavaScript - Entry point
 * 
 * Menginisialisasi komponen-komponen dasar aplikasi
 */
document.addEventListener('DOMContentLoaded', function() {
    console.log('Botopia app initialized');
    
    // Initialize Alpine.js components
    if (typeof Alpine !== 'undefined') {
        // Register Alpine.js data components for each page
        
        // Add event listener for dark mode toggle
        const darkModeToggle = document.getElementById('darkModeToggle');
        if (darkModeToggle) {
            darkModeToggle.addEventListener('click', toggleDarkMode);
        }
        
        // Initialize toast notification system
        initializeToasts();
    }
    
    // Initialize sidebar
    initializeSidebar();
});

// Initialize toasts for non-Alpine.js context
function initializeToasts() {
    window.showToast = function(type, message, duration = 3000) {
        // If the toast utility is already defined elsewhere, use that
        if (typeof window.showToast === 'function' && 
            window.showToast.toString().includes('p-toast-container')) {
            return;
        }
        
        // Otherwise implement a basic toast system
        // Get or create toast container
        let container = document.getElementById('toast-container');
        if (!container) {
            container = document.createElement('div');
            container.id = 'toast-container';
            container.className = 'p-toast-container';
            document.body.appendChild(container);
        }
        
        // Create toast element
        const toast = document.createElement('div');
        toast.className = `p-toast p-toast--${type}`;
        
        // Add icon based on type
        let icon = 'fa-info-circle';
        switch (type) {
            case 'success':
                icon = 'fa-check-circle';
                break;
            case 'error':
                icon = 'fa-circle-exclamation';
                break;
            case 'warning':
                icon = 'fa-triangle-exclamation';
                break;
        }
        
        // Set content
        toast.innerHTML = `
            <div class="p-toast__icon">
                <i class="fas ${icon}"></i>
            </div>
            <div class="p-toast__content">${message}</div>
            <button class="p-toast__close" aria-label="Close">
                <i class="fas fa-times"></i>
            </button>
        `;
        
        // Add to container
        container.appendChild(toast);
        
        // Add show class after a small delay for animation
        setTimeout(() => {
            toast.classList.add('p-toast--show');
        }, 10);
        
        // Add click event for close button
        const closeBtn = toast.querySelector('.p-toast__close');
        closeBtn.addEventListener('click', () => {
            closeToast(toast);
        });
        
        // Auto close after duration
        setTimeout(() => {
            closeToast(toast);
        }, duration);
    };
    
    function closeToast(toast) {
        // Add hide class for animation
        toast.classList.add('p-toast--hide');
        
        // Remove after animation completes
        toast.addEventListener('transitionend', () => {
            if (toast.parentNode) {
                toast.parentNode.removeChild(toast);
            }
        });
    }
}

// Initialize sidebar interactions
function initializeSidebar() {
    // Toggle sidebar collapse state
    const sidebarToggleBtn = document.getElementById('sidebarToggleBtn');
    const sidebar = document.querySelector('.l-sidebar');
    const content = document.querySelector('.l-content-with-sidebar');
    
    if (sidebarToggleBtn && sidebar && content) {
        sidebarToggleBtn.addEventListener('click', function() {
            sidebar.classList.toggle('collapsed');
            content.classList.toggle('sidebar-collapsed');
            
            // Save state in localStorage
            localStorage.setItem('sidebarCollapsed', sidebar.classList.contains('collapsed'));
        });
        
        // Restore saved state
        const isCollapsed = localStorage.getItem('sidebarCollapsed') === 'true';
        if (isCollapsed) {
            sidebar.classList.add('collapsed');
            content.classList.add('sidebar-collapsed');
        }
    }
    
    // Handle mobile sidebar toggle
    const mobileSidebarToggle = document.getElementById('mobileSidebarToggle');
    const sidebarBackdrop = document.querySelector('.l-sidebar-backdrop');
    
    if (mobileSidebarToggle && sidebar) {
        mobileSidebarToggle.addEventListener('click', function() {
            sidebar.classList.add('mobile-open');
            if (sidebarBackdrop) {
                sidebarBackdrop.classList.add('active');
            }
        });
    }
    
    // Handle backdrop click to close sidebar
    if (sidebarBackdrop) {
        sidebarBackdrop.addEventListener('click', function() {
            sidebar.classList.remove('mobile-open');
            sidebarBackdrop.classList.remove('active');
        });
    }
}

// Dark mode toggle functionality
function toggleDarkMode() {
    document.documentElement.classList.toggle('dark-mode');
    const isDarkMode = document.documentElement.classList.contains('dark-mode');
    localStorage.setItem('darkMode', isDarkMode);
}
