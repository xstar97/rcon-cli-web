const fs = require('fs');
const YAML = require('yaml');
const { spawn } = require('child_process');
const path = require('path');
const { CONFIG, COMMANDS } = require('./config');
const http = require('http');

const defaultServer = CONFIG.CLI_DEFAULT_SERVER;
const configFile = CONFIG.CLI_CONFIG;
const cliRoot = CONFIG.CLI_ROOT;
const cliOptVersion = COMMANDS.VERSION;
const cliOptEnv = COMMANDS.ENV;

function getServersFromConfig() {
    const fileContent = fs.readFileSync(configFile, 'utf8');
    const config = YAML.parse(fileContent);

    const servers = Object.keys(config).map(serverName => {
        const type = config[serverName].type || 'rcon';
        return { server: serverName, type: type };
    });

    return servers;
}

async function sendToCLI(server, command) {
    return new Promise((resolve, reject) => {
        const args = [cliOptEnv, server, command];
        const cliFile = path.join(cliRoot);
        const cliProcess = spawn(cliFile, args);

        let output = '';

        cliProcess.stdout.on('data', data => {
            output += data.toString();
        });

        cliProcess.stderr.on('data', error => {
            // Fail silently and return failure message
            resolve("Command failed to send.");
        });

        cliProcess.on('close', code => {
            if (code !== 0) {
                // Fail silently and return failure message
                resolve("Command failed to send.");
            } else {
                resolve(output);
            }
        });

        // Error event handler to catch spawn errors
        cliProcess.on('error', err => {
            // Fail silently and return failure message
            resolve("Command failed to send.");
        });
    });
}

async function checkAndUpdateVersion() {
    try {
        const options = {
            hostname: 'api.github.com',
            path: '/repos/gorcon/rcon-cli/releases/latest',
            method: 'GET',
            headers: {
                'User-Agent': 'nodejs' // GitHub requires a User-Agent header
            }
        };

        const req = http.request(options, res => {
            let data = '';
            res.on('data', chunk => {
                data += chunk;
            });
            res.on('end', () => {
                const latestVersion = JSON.parse(data).tag_name.substring(1); // Remove 'v' prefix
                sendToCLI(defaultServer, cliOptVersion)
                    .then(currentVersionResponse => {
                        const currentVersion = currentVersionResponse.split(' ')[2].trim();
                        const updateAvailable = latestVersion !== currentVersion;
                        return { latestVersion, currentVersion, updateAvailable };
                    })
                    .catch(error => {
                        console.error('Error sending command:', error);
                        return { error: 'Error sending command' };
                    });
            });
        });

        req.on('error', error => {
            console.error('Error making HTTP request:', error);
            return { error: 'Error making HTTP request' };
        });

        req.end();
    } catch (error) {
        console.error('Error checking or updating version:', error);
        return { error: 'Error checking or updating version' };
    }
}

module.exports = {
    getServersFromConfig,
    sendToCLI,
    checkAndUpdateVersion
};
