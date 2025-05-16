/**
 * Script loader dan debugging helper
 */
console.log("Botopia scripts loading...");

// Track script loading
window.BotopiaScripts = {
    loaded: [],
    failed: [],
    
    markLoaded: function(scriptName) {
        this.loaded.push(scriptName);
        console.log(`✅ Script loaded: ${scriptName}`);
    },
    
    markFailed: function(scriptName, error) {
        this.failed.push({ name: scriptName, error: error });
        console.error(`❌ Script failed: ${scriptName}`, error);
    },
    
    reportStatus: function() {
        console.log("Botopia script loading status:", {
            loaded: this.loaded,
            failed: this.failed
        });
    }
};

// Add event to check all scripts loaded
window.addEventListener('load', function() {
    console.log("Window loaded, checking script status...");
    window.BotopiaScripts.reportStatus();
    
    // Check specific Alpine components
    console.log("Alpine global components:", {
        qrPage: typeof window.qrPage === 'function' ? "defined" : "undefined",
        dashboardApp: typeof window.dashboardApp === 'function' ? "defined" : "undefined",
        statsCards: typeof window.statsCards === 'function' ? "defined" : "undefined"
    });
});

// Report script loaded
window.BotopiaScripts.markLoaded('loader-debug.js');


/**
 * Loader Debug Script - Membantu mendeteksi masalah loading script
 */
(function() {
    // Daftar file JavaScript penting
    const criticalScripts = [
        { name: 'Alpine.js', path: '/static/js/alpine.min.js' },
        { name: 'Toast.js', path: '/static/js/toast.js' },
        { name: 'Main.js', path: '/static/js/main.js' },
        { name: 'Contacts App', path: '/static/js/contacts/contacts-app.js' }
    ];
    
    // Mengecek apakah script sudah dimuat
    function checkLoadedScripts() {
        const loadedScripts = Array.from(document.scripts).map(script => script.src);
        console.log('➡️ Checking critical scripts...');
        
        criticalScripts.forEach(script => {
            const isLoaded = loadedScripts.some(src => src.includes(script.path));
            console.log(`${isLoaded ? '✅' : '❌'} ${script.name}: ${script.path}`);
            
            if (!isLoaded && script.path.includes('contacts-app.js')) {
                console.warn('⚠️ CRITICAL: Contacts App script is not loaded! This will cause button issues.');
            }
        });
        
        // Cek apakah Alpine.js terinisialisasi
        if (window.Alpine) {
            console.log('✅ Alpine.js is initialized');
        } else {
            console.warn('❌ Alpine.js is not initialized! This will cause component issues.');
        }
        
        // Cek halaman saat ini
        const currentPage = document.querySelector('body').dataset.page || 'unknown';
        console.log(`📄 Current page: ${currentPage}`);
    }
    
    // Run checks when DOM is loaded
    document.addEventListener('DOMContentLoaded', checkLoadedScripts);
})();
