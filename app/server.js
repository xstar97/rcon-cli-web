const express = require('express');
const fs = require('fs');
const { CONFIG } = require('./config');
const Keyv = require('keyv');

const app = express();

let rconCliPath = CONFIG.CLI_ROOT;

// Validate if rconCliPath exists
if (!fs.existsSync(rconCliPath)) {
    console.error(`Error: rconCliPath '${rconCliPath}' does not exist.`);
    process.exit(1); // Exit the process with an error code
}

app.use(express.json()); // Parse JSON bodies

// Determine the storage adapter based on the DB_TYPE
let storageAdapter;
if (CONFIG.DB_TYPE === "sqlite") {
    storageAdapter = `sqlite://${CONFIG.SQLITE_DB}`;
    console.log(`Using sqlite...`);
} else {
    storageAdapter = `redis://${CONFIG.REDIS_USER}:${CONFIG.REDIS_PASS}@${CONFIG.REDIS_HOST}:${CONFIG.REDIS_PORT}`;
    console.log(`Using redis...`);
}

// Create a new Keyv instance using the determined storage adapter
const keyv = new Keyv(storageAdapter);

// Handle DB connection errors
keyv.on('error', err => console.error('Keyv Connection Error', err));

// Export the Keyv instance
module.exports = keyv;

// Import route files after exporting the Keyv instance
const indexRoute = require('../routes/index');
const rconRoutes = require('../routes/rcon');
const logsRoutes = require('../routes/logs');
const savedRoutes = require('../routes/saved');

// Use route middleware
app.use('/', indexRoute); // Use the index route for serving static files
app.use('/logs', logsRoutes);
app.use('/rcon', rconRoutes);
app.use('/saved', savedRoutes);

// Start the server
const server = app.listen(CONFIG.PORT, () => {
    console.log(`Server running on port ${CONFIG.PORT}`);
});