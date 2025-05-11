/**
 * Dashboard statistics card functionality
 */
window.statsCards = function() {
    return {
        connectionState: "disconnected",
        messageCount: 0,
        commandsRun: 0,
        uptime: 0,
        
        init(state, msgCount, cmdRun) {
            this.connectionState = state || "disconnected";
            this.messageCount = parseInt(msgCount || "0");
            this.commandsRun = parseInt(cmdRun || "0");
            
            console.log("Stats cards initialized with:", {
                state: this.connectionState,
                messages: this.messageCount,
                commands: this.commandsRun
            });
            
            // Subscribe to stats updates
            document.addEventListener('stats-updated', e => {
                if (e.detail) {
                    this.connectionState = e.detail.connectionState;
                    this.messageCount = e.detail.messageCount;
                    this.commandsRun = e.detail.commandsRun;
                    this.uptime = e.detail.uptime;
                }
            });
        },
        
        formatDuration(seconds) {
            if (!seconds) return "";
            
            const days = Math.floor(seconds / 86400);
            seconds %= 86400;
            const hours = Math.floor(seconds / 3600);
            seconds %= 3600;
            const minutes = Math.floor(seconds / 60);
            
            let result = "";
            if (days > 0) result += `${days}h `;
            if (hours > 0 || days > 0) result += `${hours}j `;
            return result + `${minutes}m`;
        }
    };
};

// Report script loaded
document.addEventListener('DOMContentLoaded', function() {
    console.log("Stats cards script loaded");
});
