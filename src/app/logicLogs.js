const fs = require('fs');
const YAML = require('yaml');
const { CONFIG } = require('./config');

let configFile = CONFIG.CLI_CONFIG;
let defaultServer = CONFIG.CLI_DEFAULT_SERVER;

// Function to read the content of a log file based on the log path
function readLogFile(logPath) {
    try {
        // Check if the log file exists
        if (!fs.existsSync(logPath)) {
            return null; // Return null if the log file doesn't exist
        }

        // Read the content of the log file
        return fs.readFileSync(logPath, 'utf8');
    } catch (error) {
        console.error('Error reading log file:', error.message);
        return null;
    }
}

// Function to read rcon.yaml and extract server names and log paths
function getLogsFromConfig() {
    try {
        const fileContent = fs.readFileSync(configFile, 'utf8');
        const config = YAML.parse(fileContent);

        const logs = Object.keys(config).map(serverName => {
            const logPath = config[serverName].log || null;
            return { server: serverName, log: logPath };
        });

        return logs;
    } catch (error) {
        console.error('Error reading config file:', error.message);
        return [];
    }
}

module.exports = {
    readLogFile,
    getLogsFromConfig,
    defaultServer
};
