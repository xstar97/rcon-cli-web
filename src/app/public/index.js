document.addEventListener("DOMContentLoaded", function() {

    const commandInput = document.getElementById("commandInput");
    const submitBtn = document.getElementById("submitBtn");
    const rconOutput = document.getElementById("rconOutput");
    const toggleModeBtn = document.getElementById("toggleModeBtn");
    const body = document.body;
    const serversSelect = document.getElementById("servers");

    let currentServer = null;
    let currentMode = null;

    // Routes
    const SERVERS_ROUTE = "/rcon/servers";
    const SAVED_DATA_ROUTE = "/saved";
    const RCON_ROUTE = "/rcon";
    const VERSION_ROUTE = "/rcon/version";
    const LOGS_ROUTE = "/logs";

    // Initialize the application
    fetchServers();
    fetchSavedData();
    fetchVersion();
    
    // Function to fetch the version from the /rcon/version route
    function fetchVersion() {
        fetch(VERSION_ROUTE)
            .then(response => response.json())
            .then(data => {
                const versionElement = document.getElementById("cli-version");
                versionElement.textContent = `cli-version: ${data.currentVersion}`;
            })
            .catch(error => console.error("Error fetching version:", error));
    }
    
    // Fetch servers from the /rcon route and populate the select element
    function fetchServers() {
        fetch(SERVERS_ROUTE)
            .then(response => response.json())
            .then(data => {
                data.forEach(server => {
                    addServerOption(server.server);
                });
                currentServer = serversSelect.value;
            })
            .catch(error => console.error("Error fetching servers:", error));
    }

    // Fetch saved data from the /saved route
    function fetchSavedData() {
        fetch(SAVED_DATA_ROUTE)
            .then(response => response.json())
            .then(data => {
                currentMode = data.mode;
                currentServer = data.server;
                updateUI();
            })
            .catch(error => console.error("Error fetching saved data:", error));
    }
    
    // Function to add server options to the select element
    function addServerOption(serverName) {
        const option = document.createElement("option");
        option.value = serverName;
        option.textContent = serverName;
        serversSelect.appendChild(option);
    }

    // Function to update the UI based on saved data
    function updateUI() {
        body.classList.toggle("dark-mode", currentMode === 'dark');
        toggleModeBtn.classList.toggle("dark-mode-btn", currentMode === 'dark');
        serversSelect.value = currentServer;
    }

    // Function to clear output
    function clearOutput() {
        rconOutput.textContent = "";
    }

    // Function to send command and update logs
    function sendCommand(command) {
        fetch(RCON_ROUTE, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ server: currentServer, command })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to execute command.');
            }
            return response.json();
        })
        .then(data => {
            let content = `command: ${data.command}\n${data.output}`
            updateOutput(content);
        })
        .catch(error => console.error('Error:', error));
    }

    // Function to update output
    function updateOutput(output) {
        // Split the output string into an array of lines
        const lines = output.split("\n");
        
        // Iterate over each line and create a paragraph element for it
        lines.forEach(line => {
            const paragraph = document.createElement('p');
            paragraph.textContent = line;
            rconOutput.appendChild(paragraph);
        });
    }

    // Function to clear input
    function clearInput() {
        commandInput.value = "";
    }

    // Event listeners

    submitBtn.addEventListener("click", function() {
        const command = commandInput.value.trim();
        if (command !== "") {
            sendCommand(command);
            clearInput();
        }
    });

    commandInput.addEventListener("keypress", function(event) {
        if (event.key === "Enter") {
            const command = commandInput.value.trim();
            if (command !== "") {
                sendCommand(command);
                clearInput();
            }
        }
    });

    serversSelect.addEventListener("change", function() {
        currentServer = this.value;
        updateSavedData();
        clearOutput();
    });

    toggleModeBtn.addEventListener("click", function() {
        body.classList.toggle("dark-mode");
        toggleModeBtn.classList.toggle("dark-mode-btn");
        currentMode = body.classList.contains("dark-mode") ? "dark" : "light";
        updateSavedData();
    });

    // Function to update saved data
    function updateSavedData() {
        fetch(SAVED_DATA_ROUTE, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ server: currentServer, mode: currentMode })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to update saved data.');
            }
            return response.json();
        })
        .catch(error => console.error('Error updating saved data:', error));
    }
    
    // Open a new tab to view logs when the "View Logs" button is clicked
    const viewLogsBtn = document.getElementById("viewLogsBtn");
    viewLogsBtn.addEventListener("click", function() {
        const logsUrl = `${LOGS_ROUTE}/${currentServer}`;
        window.open(logsUrl, '_blank');
    });
});
