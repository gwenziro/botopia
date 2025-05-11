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
