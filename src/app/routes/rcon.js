const express = require('express');
const router = express.Router();
const { getServersFromConfig, sendToCLI, checkAndUpdateVersion } = require('../logic/rcon.js');

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
        res.status(500).json({ error: 'Error executing command' });
    }
});

// Version route
router.get('/version', async (req, res) => {
    try {
        const versionInfo = await checkAndUpdateVersion();
        res.json(versionInfo);
    } catch (error) {
        console.error('Error getting version:', error);
        res.status(500).json({ error: 'Internal server error' });
    }
});

module.exports = router;
