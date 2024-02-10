const keyv = require('./server.js'); // Import the Keyv instance directly
const { CONFIG } = require('./config.js');

const defaultServer = CONFIG.CLI_DEFAULT_SERVER;
const defaultMode = CONFIG.MODE;

async function getSavedData() {
    try {
        const saved = {
            server: await keyv.get('server') || defaultServer,
            mode: await keyv.get('mode') || defaultMode
        };
        return saved;
    } catch (error) {
        console.error('Error retrieving saved data:', error);
        throw new Error('Error retrieving saved data');
    }
}

async function getSavedValueByKey(key) {
    try {
        const value = await keyv.get(key);
        if (value === undefined) {
            throw new Error('Key not found');
        } else {
            return { [key]: value };
        }
    } catch (error) {
        console.error('Error retrieving key:', error);
        throw new Error('Error retrieving key');
    }
}

async function saveData(server, mode) {
    try {
        // Update server and mode separately
        if (server !== undefined) {
            await keyv.set('server', server);
        }
        if (mode !== undefined) {
            await keyv.set('mode', mode);
        }
        return { success: true };
    } catch (error) {
        console.error('Error saving data:', error);
        throw new Error('Error saving data');
    }
}

module.exports = {
    getSavedData,
    getSavedValueByKey,
    saveData
};
