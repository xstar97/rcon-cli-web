const express = require('express');
const router = express.Router();
const keyv = require('./server'); // Import the Keyv instance directly

const { CONFIG } = require('./config');

// Define saved routes
router.get('/', async (req, res) => {
    try {
        const saved = {
            server: await keyv.get('server') || CONFIG.CLI_DEFAULT_SERVER,
            mode: await keyv.get('mode') || CONFIG.MODE
        };
        res.json(saved);
    } catch (error) {
        console.error('Error retrieving saved data:', error);
        res.status(500).json({ error: 'Error retrieving saved data' });
    }
});

router.get('/:key', async (req, res) => {
    try {
        const { key } = req.params;
        const value = await keyv.get(key);
        if (value === undefined) {
            res.status(404).json({ error: 'Key not found' });
        } else {
            res.json({ [key]: value });
        }
    } catch (error) {
        console.error('Error retrieving key:', error);
        res.status(500).json({ error: 'Error retrieving key' });
    }
});

router.post('/', async (req, res) => {
    try {
        const { server, mode } = req.body;

        // Update server and mode separately
        if (server !== undefined) {
            await keyv.set('server', server);
        }
        if (mode !== undefined) {
            await keyv.set('mode', mode);
        }
        res.json({ success: true });
    } catch (error) {
        console.error('Error saving data:', error);
        res.status(500).json({ error: 'Error saving data' });
    }
});

module.exports = router;
