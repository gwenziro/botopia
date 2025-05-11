/**
 * Main dashboard application logic
 */

// Define the Alpine component function globally
window.dashboardApp = function() {
    return {
        connectionState: "disconnected",
        commandCount: 0,
        messageCount: 0,
        commandsRun: 0,
        uptime: 0,
        commands: {}, 
        initialized: false,
        isRefreshing: false,
        refreshTimer: null,
        phone: "",  // Add phone field
        name: "",   // Add name field

        init() {
            // Load initial values from HTML data attributes
            const dashboardEl = document.getElementById('dashboard-app');
            if (dashboardEl) {
                this.connectionState = dashboardEl.dataset.connectionState || "disconnected";
                this.commandCount = parseInt(dashboardEl.dataset.commandCount || "0");
                this.messageCount = parseInt(dashboardEl.dataset.messageCount || "0");
                this.commandsRun = parseInt(dashboardEl.dataset.commandsRun || "0");
                this.phone = dashboardEl.dataset.phone || "";
                this.name = dashboardEl.dataset.name || "";
            }
            
            // Load commands
            this.loadCommands();
            
            // Start polling for updates
            this.fetchLatestStats();
            this.refreshTimer = setInterval(() => this.fetchLatestStats(), 5000);
            
            this.initialized = true;
        },
        
        loadCommands() {
            try {
                const commandsData = document.getElementById('commands-data');
                if (!commandsData) {
                    console.warn("Commands data element not found");
                    return;
                }
                
                const rawData = commandsData.textContent.trim();
                if (!rawData) {
                    console.warn("Commands data is empty");
                    return;
                }
                
                let parsed;
                try {
                    parsed = JSON.parse(rawData);
                } catch (e) {
                    console.error("Failed to parse JSON:", e);
                    return;
                }
                
                // Use helper function for consistent formatting
                this.commands = getFormattedCommands(parsed);
                console.log(`Loaded ${Object.keys(this.commands).length} commands:`, 
                    Object.keys(this.commands));
            } catch (error) {
                console.error('Error in loadCommands:', error);
            }
        },

        refreshStats() {
            if (this.isRefreshing) return;
            
            this.isRefreshing = true;
            
            // Call fetchLatestStats and ensure isRefreshing is reset after completion
            this.fetchLatestStats(true)
                .catch(error => console.error("Error during refresh:", error))
                .finally(() => {
                    // Reset isRefreshing after a short delay for better UX
                    setTimeout(() => {
                        this.isRefreshing = false;
                    }, 800);
                });
        },

        fetchLatestStats(showNotification = false) {
            // Return the promise so we can chain .finally() to it
            return fetch('/api/stats')
                .then(response => {
                    if (!response.ok) {
                        throw new Error(`HTTP error! Status: ${response.status}`);
                    }
                    return response.json();
                })
                .then(data => {
                    // Check if connection state changed
                    const wasConnected = this.connectionState === 'connected';
                    const isNowConnected = data.connectionState === 'connected';
                    
                    // Update data
                    this.connectionState = data.connectionState;
                    this.messageCount = data.messageCount;
                    this.commandsRun = data.commandsRun;
                    this.uptime = data.uptime;
                    
                    // Update contact info
                    if (data.phone) this.phone = data.phone;
                    if (data.name) this.name = data.name;
                    
                    // Show notification if connection state changed
                    if (wasConnected !== isNowConnected) {
                        window.showToast({
                            title: isNowConnected ? 'Connected' : 'Disconnected',
                            message: isNowConnected 
                                ? 'Bot is now connected to WhatsApp' 
                                : 'Bot has been disconnected from WhatsApp',
                            type: isNowConnected ? 'success' : 'error'
                        });
                    }
                    
                    // Optional notification when manually refreshed
                    if (showNotification) {
                        window.showToast({
                            title: 'Stats Updated',
                            message: 'Dashboard statistics have been refreshed',
                            type: 'info',
                            duration: 3000
                        });
                    }
                })
                .catch(error => {
                    console.error('Error fetching stats:', error);
                    
                    if (showNotification) {
                        window.showToast({
                            title: 'Error',
                            message: 'Failed to update dashboard statistics',
                            type: 'error'
                        });
                    }
                    
                    // Re-throw to propagate to caller
                    throw error;
                });
        },

        formatDuration(seconds) {
            if (!seconds) return "Just now";
            
            const days = Math.floor(seconds / 86400);
            seconds %= 86400;
            const hours = Math.floor(seconds / 3600);
            seconds %= 3600;
            const minutes = Math.floor(seconds / 60);
            const remainingSeconds = Math.floor(seconds % 60);
            
            let result = "";
            
            if (days > 0) {
                result += `${days}d `;
            }
            
            if (hours > 0 || days > 0) {
                result += `${hours}h `;
            }
            
            if (minutes > 0 || hours > 0 || days > 0) {
                result += `${minutes}m `;
            }
            
            // Only show seconds if less than 1 hour total duration
            if ((days === 0 && hours === 0) || remainingSeconds > 0) {
                result += `${remainingSeconds}s`;
            }
            
            return result.trim();
        },

        // Helper functions for displaying contact info
        getDisplayPhone() {
            return this.phone || "Unknown number";
        },
        
        getDisplayName() {
            return "WhatsApp User"; // Selalu tampilkan nilai default
        }
    };
};

// Additional function to ensure proper command processing
function getFormattedCommands(commandsData) {
    // Ensure data exists
    if (!commandsData) {
        return {};
    }
    
    // Log for debugging
    console.log("Raw commands data:", typeof commandsData, commandsData);
    
    // If already an object, return directly
    if (typeof commandsData === 'object' && commandsData !== null && !Array.isArray(commandsData)) {
        return commandsData;
    }
    
    // If string, parse JSON
    if (typeof commandsData === 'string') {
        try {
            return JSON.parse(commandsData);
        } catch (e) {
            console.error("Failed to parse commands JSON:", e);
            return {};
        }
    }
    
    // If array, convert to object
    if (Array.isArray(commandsData)) {
        console.warn("Commands data is an array, converting to object");
        const result = {};
        commandsData.forEach((cmd, index) => {
            if (cmd && typeof cmd === 'object' && cmd.name) {
                result[cmd.name] = cmd;
            } else {
                result[`command-${index}`] = cmd;
            }
        });
        return result;
    }
    
    // Fallback
    console.error("Unexpected command data format:", commandsData);
    return {};
}

// Report script loaded
document.addEventListener('DOMContentLoaded', function() {
    console.log("Dashboard app script loaded");
});
