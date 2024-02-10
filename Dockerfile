# Use a lightweight Node.js image as base
FROM node:21.6.1-alpine3.19@sha256:a4846d0aac29ceb77a633edcbc56260231fe6f0ba3aeca1ed8f3085a26c54f8e

# Set the working directory inside the container
WORKDIR /src/app

# Install curl and tar
RUN apk add --no-cache curl tar

# Create directories for volumes
RUN mkdir -p /config /src/app/rcon

# Set permissions for the /rcon and /config directories
RUN chmod -R 755 /src/app/rcon /config

# Download the latest release of rcon-cli
RUN curl -L -o /tmp/rcon.tar.gz $(curl -s https://api.github.com/repos/gorcon/rcon-cli/releases/latest | grep "browser_download_url.*amd64_linux.tar.gz" | cut -d '"' -f 4)

# Extract rcon binary and rcon.yaml configuration file
RUN tar -xzf /tmp/rcon.tar.gz -C /src/app/rcon && \
    mv /src/app/rcon/rcon-*/* /src/app/rcon && \
    rm -rf /src/app/rcon/rcon-*

# Copy the rest of the application code from app directory to /src/app
COPY src/app/* /src/app/

# Install dependencies
RUN npm install --production

# Validate that files were copied successfully
RUN ls -al /src/app/

# Set environment variables
ENV PORT=3000 \
    NODE_ENV=production \
    MODE=dark \
    CLI_ROOT=/src/app/rcon \
    CLI_CONFIG=/src/app/rcon/rcon.yaml \
    CLI_DEFAULT_SERVER=default \
    DB_TYPE=sqlite \
    SQLITE_DB=/config/rcon.sqlite

# Expose the port defined in the environment variable
EXPOSE $PORT

# Arguments to specify user and group names with default values
ARG USER=kah
ARG GROUP=kah

# Create a new user and group
RUN addgroup -S $GROUP && adduser -S $USER -G $GROUP

# Set permissions for the /config directory
RUN chown -R $USER:$GROUP /config

# Switch to the newly created user
USER $USER

# Define volumes
VOLUME ["/config"]

# Run the Node.js application using `npm start` script
CMD ["npm", "start"]

# Health check
HEALTHCHECK --interval=60s --timeout=10s --start-period=5s --retries=6 \
    CMD curl -fs http://localhost:$PORT/ || exit 1
