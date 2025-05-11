/**
 * Dashboard command list functionality
 */
window.commandList = function() {
    return {
        commands: {},
        
        init(commandsData) {
            try {
                this.commands = JSON.parse(commandsData);
                this.renderCommands();
            } catch(e) {
                console.error('Failed to parse commands:', e);
            }
        },
        
        renderCommands() {
            const tableBody = document.getElementById('commands-list');
            if (!tableBody) return;
            
            tableBody.innerHTML = '';
            
            // Sort command names alphabetically
            const commandNames = Object.keys(this.commands).sort();
            
            for (const name of commandNames) {
                const cmd = this.commands[name];
                const row = document.createElement('tr');
                row.className = 'hover:bg-gray-50';
                
                row.innerHTML = `
                    <td class="px-6 py-4 whitespace-nowrap">
                        <div class="flex items-center">
                            <div class="text-sm font-medium text-gray-900">
                                <span class="font-mono bg-gray-100 py-1 px-2 rounded">!${name}</span>
                            </div>
                        </div>
                    </td>
                    <td class="px-6 py-4 text-sm text-gray-500">${cmd.Description || ''}</td>
                `;
                
                tableBody.appendChild(row);
            }
            
            console.log(`Command list rendered with ${commandNames.length} commands`);
        }
    };
};

// Report script loaded
document.addEventListener('DOMContentLoaded', function() {
    console.log("Command list script loaded");
});
