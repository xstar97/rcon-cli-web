const fs = require('fs');
const YAML = require('yaml');
const { spawn } = require('child_process');
const path = require('path');
const { CONFIG, COMMANDS } = require('../config');

const defaultServer = CONFIG.CLI_DEFAULT_SERVER;
const configFile = CONFIG.CLI_CONFIG;

const cliRoot = CONFIG.CLI_ROOT;
const cliConfig = CONFIG.CLI_CONFIG;

const cmdVER = COMMANDS.VERSION;
const cmdENV = COMMANDS.ENV;
const cmdConfig = COMMANDS.CONFIG

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
    console.log(`Sending command to CLI:\nServer: ${server}\nCommand: ${command}`);

    return new Promise((resolve, reject) => {
        const args = [cmdConfig,cliConfig,cmdENV, server, command];

        const cliFile = path.join(cliRoot);

        const cliProcess = spawn(cliFile, args);

        let output = '';
        let errorOutput = '';

        cliProcess.stdout.on('data', data => {
            output += data.toString();
        });

        cliProcess.stderr.on('data', data => {
            errorOutput += data.toString();
            console.error(`CLI stderr: ${data.toString()}`);
        });

        cliProcess.on('close', code => {
            if (code !== 0) {
                console.error(`CLI process exited with code ${code}`);
                console.error(`Error output: ${errorOutput}`);
                resolve("Command failed to send.");
            } else {
                console.log(`CLI process exited successfully`);
                resolve(output);
            }
        });

        cliProcess.on('error', err => {
            console.error(`CLI process encountered an error: ${err.message}`);
            resolve("Command failed to send.");
        });
    });
}

async function checkAndUpdateVersion() {
    try {
        const response = await fetch('https://api.github.com/repos/gorcon/rcon-cli/releases/latest');
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        const data = await response.json();
        const latestVersion = data.tag_name.substring(1); // Remove 'v' prefix

        const currentVersionResponse = await sendToCLI(defaultServer, cmdVER);
        const currentVersion = currentVersionResponse.split(' ')[2].trim();

        const updateAvailable = latestVersion !== currentVersion;

        return { latestVersion, currentVersion, updateAvailable };
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
