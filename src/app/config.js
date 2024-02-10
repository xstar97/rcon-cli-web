const path = require('path');
const dotenv = require('dotenv');

// Set NODE_ENV to 'development' if not already set
process.env.NODE_ENV = process.env.NODE_ENV || 'development';

// Load different environment variables based on NODE_ENV
if (process.env.NODE_ENV === 'production') {
    dotenv.config();
} else {
    dotenv.config({ path: path.resolve(__dirname, '../../development.env') });
}

// Configuration constants
const CONFIG = {
    //web port
    PORT: process.env.PORT || 3000,
    //dark/light
    MODE: process.env.MODE || "dark",
    //Root path to rcon file
    CLI_ROOT: path.join(process.env.CLI_ROOT || "/home/kah/app/rcon/rcon"),
    //Root path to rcon.yaml
    CLI_CONFIG: path.join(process.env.CLI_CONFIG || "/config/rcon.yaml"),
    //Default rcon env
    CLI_DEFAULT_SERVER: process.env.CLI_DEFAULT_SERVER || "default",
    //sqlite | redis
    DB_TYPE: process.env.DB_TYPE || "sqlite",
    //redis host
    REDIS_HOST: process.env.REDIS_HOST, 
    //redis port
    REDIS_PORT: process.env.REDIS_PORT || 6379,
    //redis user
    REDIS_USER: process.env.REDIS_USER || "default",
    //redis pass
    REDIS_PASS: process.env.REDIS_PASS,
    //sqlite
    SQLITE_DB:`${path.join(process.env.SQLITE_DB || "/config/db/sqlite.db")}`
};

const COMMANDS = {
    VERSION: "--version",
    ENV: "--env",
    CONFIG: "--config" 
}

module.exports = { CONFIG, COMMANDS };
