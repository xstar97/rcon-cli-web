const express = require('express');
const router = express.Router();
const fs = require('fs');
const YAML = require('yaml');
const { CONFIG, COMMANDS } = require('../config');
const { spawn } = require('child_process');

const defaultServer = CONFIG.CLI_DEFAULT_SERVER;
const configFile = CONFIG.CLI_CONFIG;
const cliRoot = CONFIG.CLI_ROOT;

// Function to read rcon.yaml and extract server names
function getServersFromConfig() {
    const fileContent = fs.readFileSync(configFile, 'utf8');
    const config = YAML.parse(fileContent);

    const servers = Object.keys(config).map(serverName => {
        const type = config[serverName].type || 'rcon';
        return { server: serverName, type: type };
    });

    return servers;
}

// Servers route
router.get('/servers', (req, res) => {
    try {
        const servers = getServersFromConfig();
        res.json(servers);
    } catch (error) {
        console.error('Error getting server names:', error);
        res.status(500).json({ error: 'Internal server error' });
    }
});

// RCON route
router.post('/', async (req, res) => {
    const { server, command } = req.body;
    const selectedServer = server || defaultServer;

    if (!command) {
        return res.status(400).json({ error: 'Command parameter is missing.' });
    }

    try {
        const output = await sendToCLI(selectedServer, command);
        res.json({ server: selectedServer, command, output });
    } catch (error) {
        console.error('Error executing command:', error);
        //res.status(500).json({ server: selectedServer, command, output: 'Error executing command' });
        res.json({ server: selectedServer, command, output: error });
    }
});

// Version route
router.get('/version', async (req, res) => {
    try {
        const versionOutput = await sendToCLI(CONFIG.CLI_DEFAULT_SERVER, COMMANDS.VERSION);
        const version = versionOutput.split(' ')[2].trim();
        res.json({ version });
    } catch (error) {
        console.error('Error getting version:', error);
        res.status(500).json({ error: 'Internal server error' });
    }
});

function sendToCLI(server, command) {
    return new Promise((resolve, reject) => {
        const args = [COMMANDS.ENV, server, command];
        const cliProcess = spawn(cliRoot, args);

        let output = '';

        cliProcess.stdout.on('data', data => {
            output += data.toString();
        });

        cliProcess.stderr.on('data', error => {
            reject(error.toString());
        });

        cliProcess.on('close', code => {
            if (code !== 0) {
                reject(`CLI tool exited with code ${code}`);
            } else {
                resolve(output);
            }
        });
    });
}

module.exports = router;
