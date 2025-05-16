// Helper untuk debugging Alpine.js
window.addEventListener('DOMContentLoaded', () => {
    console.log('DOM loaded, ready for debugging');
    
    // Debugging untuk Alpine.js
    document.addEventListener('alpine:init', () => {
        console.log('Alpine initialized');
        
        // Debug mode for Alpine
        window.Alpine.debug = true;
    });
    
    // Custom event tracking
    document.addEventListener('click', (e) => {
        // Debug tombol dengan @click handler
        if (e.target && e.target.hasAttribute && e.target.hasAttribute('@click')) {
            console.log('Click on element with @click attribute:', e.target);
        }
        
        // Debug tombol Tambah Kontak
        if (e.target && e.target.textContent && e.target.textContent.includes('Tambah Kontak')) {
            console.log('Click on "Tambah Kontak" button:', e.target);
        }
    });
});
