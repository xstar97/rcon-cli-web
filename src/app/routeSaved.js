const express = require('express');
const router = express.Router();
const { getSavedData, getSavedValueByKey, saveData } = require('./logicSaved.js');

// Define saved routes
router.get('/', async (req, res) => {
    try {
        const saved = await getSavedData();
        res.json(saved);
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

router.get('/:key', async (req, res) => {
    try {
        const { key } = req.params;
        const value = await getSavedValueByKey(key);
        res.json(value);
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

router.post('/', async (req, res) => {
    try {
        const { server, mode } = req.body;
        const result = await saveData(server, mode);
        res.json(result);
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

module.exports = router;
