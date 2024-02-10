const express = require('express');
const router = express.Router();
const { readLogFile, getLogsFromConfig,defaultServer } = require('../logic/logicLogs');

// Define logs route
router.get('/:server', async (req, res) => {
    try {
        const { server } = req.params;

        // Find the log path for the requested server
        const logs = getLogsFromConfig();
        const logInfo = logs.find(log => log.server === server);

        if (!logInfo) {
            res.status(404).json({ error: 'Log not found' });
            return;
        }

        // Read the content of the log file using the log path
        const logContent = readLogFile(logInfo.log);

        if (!logContent) {
            res.status(404).json({ error: 'Log not found' });
            return;
        }

        // Set the Content-Type header to text/plain
        res.setHeader('Content-Type', 'text/plain');

        // Send the log content as response
        res.send(logContent);
    } catch (error) {
        console.error('Error retrieving log:', error);
        res.status(500).json({ error: 'Error retrieving log' });
    }
});

// Handle requests to /logs endpoint
router.get('/', (req, res) => {
    try {
        // Read the content of the default log file
        const logs = getLogsFromConfig();
        const defaultLogInfo = logs.find(log => log.server === defaultServer);

        if (!defaultLogInfo) {
            res.status(404).json({ error: 'Default log not found' });
            return;
        }

        // Read the content of the default log file using the log path
        const defaultLogContent = readLogFile(defaultLogInfo.log);

        if (!defaultLogContent) {
            res.status(404).json({ error: 'Default log not found' });
            return;
        }

        // Set the Content-Type header to text/plain
        res.setHeader('Content-Type', 'text/plain');

        // Send the default log content as response
        res.send(defaultLogContent);
    } catch (error) {
        console.error('Error retrieving default log:', error);
        res.status(500).json({ error: 'Error retrieving default log' });
    }
});

module.exports = router;
