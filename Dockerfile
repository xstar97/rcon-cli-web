# Use a lightweight Node.js image as base
FROM node:20.11.0-alpine3.19@sha256:9b61ed13fef9ca689326f40c0c0b4da70e37a18712f200b4c66d3b44fd59d98e

# Set the working directory inside the container
WORKDIR /home/node/app

# Set environment variables
ENV PORT=3000 \
    NODE_ENV=production \
    MODE=dark \
    CLI_ROOT=/home/node/app/rcon/rcon \
    CLI_CONFIG=/config/rcon.yaml \
    CLI_DEFAULT_SERVER=default \
    DB_TYPE=sqlite \
    SQLITE_DB=/config/db/rcon.sqlite

# Install curl and tar
RUN apk add --no-cache curl tar

# Create directories if they don't exist and set permissions
RUN mkdir -p /config /home/node/app/rcon

# Download the latest release of rcon-cli
RUN curl -L -o /tmp/rcon.tar.gz $(curl -s https://api.github.com/repos/gorcon/rcon-cli/releases/latest | grep "browser_download_url.*amd64_linux.tar.gz" | cut -d '"' -f 4)

# Extract rcon binary and rcon.yaml configuration file
RUN tar -xzf /tmp/rcon.tar.gz -C /home/node/app/rcon --strip-components=1 && \
    rm /tmp/rcon.tar.gz

# Set permissions for certain directories
RUN chown -R node:node /home/node/app /config

# Switch to the non-root user node
USER node

# Copy the application code with preserving directories
COPY --chown=node:node src/app/public /home/node/app/public
COPY --chown=node:node src/app/routes /home/node/app/routes
COPY --chown=node:node src/app/logic /home/node/app/logic
COPY --chown=node:node src/app/*.js /home/node/app/
COPY --chown=node:node src/app/package*.json /home/node/app/

# Install dependencies
RUN npm install --production

# Expose the port defined in the environment variable
EXPOSE $PORT

# Define volumes
VOLUME ["/config"]

# Run the Node.js application using `npm start` script
CMD ["npm", "start"]

# Health check
HEALTHCHECK --interval=60s --timeout=10s --start-period=5s --retries=6 \
    CMD curl -fs http://localhost:$PORT/ || exit 1
