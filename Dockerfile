# Use a lightweight Node.js image as base
FROM node:20.11.0-alpine3.19@sha256:9b61ed13fef9ca689326f40c0c0b4da70e37a18712f200b4c66d3b44fd59d98e

# Set the working directory inside the container
WORKDIR /home/kah/app

# Set environment variables
ENV PORT=3000 \
    NODE_ENV=production \
    MODE=dark \
    CLI_ROOT=/home/kah/app/rcon/rcon \
    CLI_CONFIG=/config/rcon.yaml \
    CLI_DEFAULT_SERVER=default \
    DB_TYPE=sqlite \
    SQLITE_DB=/config/db/rcon.sqlite

# Install curl and tar
RUN apk add --no-cache curl tar

# Copy the application code with preserving directories
COPY src/app/public /home/kah/app/public
COPY src/app/routes /home/kah/app/routes
COPY src/app/logic /home/kah/app/logic
COPY src/app/*.js /home/kah/app/
COPY src/app/package*.json /home/kah/app/

# Create directories if they dont exist.
RUN mkdir -p /config /home/kah/app/rcon

# Download the latest release of rcon-cli
RUN curl -L -o /tmp/rcon.tar.gz $(curl -s https://api.github.com/repos/gorcon/rcon-cli/releases/latest | grep "browser_download_url.*amd64_linux.tar.gz" | cut -d '"' -f 4)

# Extract rcon binary and rcon.yaml configuration file
RUN tar -xzf /tmp/rcon.tar.gz -C /home/kah/app/rcon --strip-components=1 && \
    rm /tmp/rcon.tar.gz

# Install dependencies
RUN npm install --production

# Expose the port defined in the environment variable
EXPOSE $PORT

# Create a new user and group
RUN addgroup -S kah && adduser -S kah -G kah

# Set permissions for the .npm directory
RUN mkdir -p /.npm && \
    chown -R root:root /.npm

# Switch to the newly created user
USER kah

# Define volumes
VOLUME ["/config"]

# Run the Node.js application using `npm start` script
CMD ["npm", "start"]

# Health check
HEALTHCHECK --interval=60s --timeout=10s --start-period=5s --retries=6 \
    CMD curl -fs http://localhost:$PORT/ || exit 1
