# Use a lightweight Node.js image as base
FROM node:alpine3.19@sha256:a4846d0aac29ceb77a633edcbc56260231fe6f0ba3aeca1ed8f3085a26c54f8e

# Set the working directory inside the container
WORKDIR /config

# Install curl and tar
RUN apk add --no-cache curl tar

# Create directories for volumes
RUN mkdir -p /config /rcon

# Download the latest release of rcon-cli
RUN curl -L -o /tmp/rcon.tar.gz $(curl -s https://api.github.com/repos/gorcon/rcon-cli/releases/latest | grep "browser_download_url.*amd64_linux.tar.gz" | cut -d '"' -f 4)

# Extract rcon binary and rcon.yaml configuration file
RUN tar -xzf /tmp/rcon.tar.gz -C /tmp && \
    mv /tmp/rcon-*-amd64_linux/rcon /rcon && \
    mv /tmp/rcon-*-amd64_linux/rcon.yaml /rcon && \
    rm -rf /tmp/rcon-*

# Copy package.json and package-lock.json to install dependencies
COPY package*.json ./

# Install dependencies
RUN npm install --production

# Copy the rest of the application code from app directory to /config
COPY app .

# Set environment variable for port (default to 3000 if not provided)
ENV PORT=3000
ENV NODE_ENV=production

# Expose the port defined in the environment variable
EXPOSE $PORT

# Arguments to specify user and group names with default values
ARG USER=kah
ARG GROUP=kah

# Create a new user and group
RUN addgroup -S $GROUP && adduser -S $USER -G $GROUP

# Set permissions for the /config directory
RUN chown -R $USER:$GROUP /config
RUN chmod -R 755 /config

# Switch to the newly created user
USER $USER

# Define volumes
VOLUME ["/config", "/rcon"]

# Run the Node.js application using `npm start` script
CMD ["npm", "start"]