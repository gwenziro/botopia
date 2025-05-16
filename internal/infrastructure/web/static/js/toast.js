/**
 * Toast notification system for Botopia
 */

// Create toast container if it doesn't exist
function ensureToastContainer() {
  let container = document.querySelector('.toast-container');
  if (!container) {
    container = document.createElement('div');
    container.className = 'toast-container';
    document.body.appendChild(container);
  }
  return container;
}

// Create a new toast notification
function createToast(type, message, options = {}) {
  const container = ensureToastContainer();
  
  // Default options
  const defaults = {
    title: getDefaultTitle(type),
    duration: 5000,
    closable: true,
    icon: getIconForType(type)
  };
  
  // Merge options
  const settings = { ...defaults, ...options };
  
  // Create toast element
  const toast = document.createElement('div');
  toast.className = `toast toast-${type}`;
  toast.innerHTML = `
    <div class="toast-icon">
      <i class="${settings.icon}"></i>
    </div>
    <div class="toast-content">
      <div class="toast-title">${settings.title}</div>
      <div class="toast-message">${message}</div>
    </div>
    ${settings.closable ? '<button class="toast-close">&times;</button>' : ''}
  `;
  
  // Add to container
  container.appendChild(toast);
  
  // Add close functionality
  if (settings.closable) {
    const closeBtn = toast.querySelector('.toast-close');
    closeBtn.addEventListener('click', () => removeToast(toast));
  }
  
  // Auto dismiss
  if (settings.duration) {
    setTimeout(() => removeToast(toast), settings.duration);
  }
  
  return toast;
}

// Remove a toast with animation
function removeToast(toast) {
  toast.style.animation = 'toast-out-bottom 0.3s forwards';
  setTimeout(() => {
    if (toast.parentNode) {
      toast.parentNode.removeChild(toast);
    }
    
    // Remove container if empty
    const container = document.querySelector('.toast-container');
    if (container && container.children.length === 0) {
      document.body.removeChild(container);
    }
  }, 300);
}

// Get default title based on type
function getDefaultTitle(type) {
  switch (type) {
    case 'success': return 'Berhasil';
    case 'error': return 'Error';
    case 'info': return 'Informasi';
    case 'warning': return 'Perhatian';
    default: return 'Notifikasi';
  }
}

// Get icon based on type
function getIconForType(type) {
  switch (type) {
    case 'success': return 'fas fa-check-circle';
    case 'error': return 'fas fa-exclamation-circle';
    case 'info': return 'fas fa-info-circle';
    case 'warning': return 'fas fa-exclamation-triangle';
    default: return 'fas fa-bell';
  }
}

// Public API
window.showToast = function(type, message, options) {
  return createToast(type, message, options);
};

window.successToast = function(message, options) {
  return createToast('success', message, options);
};

window.errorToast = function(message, options) {
  return createToast('error', message, options);
};

window.infoToast = function(message, options) {
  return createToast('info', message, options);
};

window.warningToast = function(message, options) {
  return createToast('warning', message, options);
};
