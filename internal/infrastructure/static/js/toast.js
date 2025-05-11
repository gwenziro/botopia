/**
 * Sistema notifikasi toast
 */
window.toastSystem = function() {
    return {
        toasts: [],
        nextId: 1,
        soundEnabled: true, // Default setting with sound enabled
        
        init() {
            // Load sound preference from localStorage
            const savedPreference = localStorage.getItem('botopia_notification_sound');
            if (savedPreference !== null) {
                this.soundEnabled = savedPreference === 'true';
            }
            
            // Add a global access point to the sound toggle
            window.toggleNotificationSound = () => {
                this.soundEnabled = !this.soundEnabled;
                localStorage.setItem('botopia_notification_sound', this.soundEnabled);
                
                // Show feedback toast
                window.showToast({
                    title: 'Notification Sound',
                    message: this.soundEnabled ? 'Sound notifications enabled' : 'Sound notifications disabled',
                    type: 'info',
                    duration: 3000
                });
                
                return this.soundEnabled;
            };
            
            window.isNotificationSoundEnabled = () => this.soundEnabled;
        },
        
        add(toast) {
            const id = this.nextId++;
            const newToast = {
                id: id,
                title: toast.title || '',
                message: toast.message || '',
                type: toast.type || 'info',
                duration: toast.duration || 5000,
                progress: 100,
                visible: true
            };
            
            this.toasts.push(newToast);
            
            // Set up automatic progress reduction
            const startTime = Date.now();
            const interval = setInterval(() => {
                const elapsedTime = Date.now() - startTime;
                const remainingPercentage = Math.max(0, 100 - (elapsedTime / newToast.duration) * 100);
                
                // Find the toast and update progress
                const index = this.toasts.findIndex(t => t.id === id);
                if (index !== -1) {
                    this.toasts[index].progress = remainingPercentage;
                }
                
                if (remainingPercentage <= 0) {
                    clearInterval(interval);
                    this.remove(id);
                }
            }, 100);
            
            // Only play sound if enabled
            if (this.soundEnabled) {
                this.playSound(newToast.type);
            }
            
            return id;
        },
        
        remove(id) {
            const index = this.toasts.findIndex(toast => toast.id === id);
            if (index !== -1) {
                // Mark toast as invisible first (for animation)
                this.toasts[index].visible = false;
                
                // Remove after animation completes
                setTimeout(() => {
                    this.toasts = this.toasts.filter(toast => toast.id !== id);
                }, 300);
            }
        },
        
        playSound(type) {
            // Optional: Add sound effect based on notification type
            if (!window.AudioContext && !window.webkitAudioContext) return;
            
            const AudioContext = window.AudioContext || window.webkitAudioContext;
            const audioCtx = new AudioContext();
            
            // Create oscillator
            const oscillator = audioCtx.createOscillator();
            const gainNode = audioCtx.createGain();
            
            // Set tone based on notification type
            switch(type) {
                case 'success':
                    oscillator.frequency.value = 800;
                    break;
                case 'error':
                    oscillator.frequency.value = 300;
                    break;
                case 'warning':
                    oscillator.frequency.value = 500;
                    break;
                case 'info':
                default:
                    oscillator.frequency.value = 600;
                    break;
            }
            
            gainNode.gain.value = 0.1;
            oscillator.connect(gainNode);
            gainNode.connect(audioCtx.destination);
            
            oscillator.start();
            gainNode.gain.exponentialRampToValueAtTime(0.001, audioCtx.currentTime + 0.5);
            setTimeout(() => {
                oscillator.stop();
            }, 500);
        }
    };
};

// Make the sound toggle functionality globally accessible
window.toggleNotificationSound = function() {
    // This will be overwritten by the actual function in the toastSystem
    console.warn("Toast system not initialized yet");
    return false;
};

window.isNotificationSoundEnabled = function() {
    // This will be overwritten by the actual function in the toastSystem
    return true;
};

/**
 * Function untuk menampilkan toast dari mana saja
 */
window.showToast = function(options) {
    window.dispatchEvent(new CustomEvent('notify', { 
        detail: options 
    }));
};

// Document ready
document.addEventListener('DOMContentLoaded', function() {
    console.log('Toast system initialized');
});
