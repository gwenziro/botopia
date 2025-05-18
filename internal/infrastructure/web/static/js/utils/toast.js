/**
 * Toast Notification System
 * Menampilkan notifikasi toast
 */
const showToast = (type, message, duration = 3000) => {
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

const closeToast = (toast) => {
    // Add hide class for animation
    toast.classList.add('p-toast--hide');
    
    // Remove after animation completes
    toast.addEventListener('transitionend', () => {
        if (toast.parentNode) {
            toast.parentNode.removeChild(toast);
        }
    });
};
