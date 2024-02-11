# Stage 1 - Build the Go application and download CLI
FROM golang:1.19.3-alpine AS builder

# Install necessary build dependencies
RUN apk --no-cache add --update gcc musl-dev curl tar

# Create the necessary directories
WORKDIR /build

# Copy Go module files
COPY cmd/go.mod cmd/go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY cmd/* ./

# Build the Go application
ARG VERSION=docker
RUN CGO_ENABLED=1 go build -o /output/rcon-cli-web .

# Stage 2 - Create the final image
FROM alpine AS runner

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /output/rcon-cli-web /app/

# Copy the public directory
COPY public/* /app/public/

# Download the latest release of rcon-cli
RUN mkdir -p /app/rcon \
    && curl -L -o /tmp/rcon.tar.gz $(curl -s https://api.github.com/repos/gorcon/rcon-cli/releases/latest | grep "browser_download_url.*amd64_linux.tar.gz" | cut -d '"' -f 4) \
    && tar -xzf /tmp/rcon.tar.gz -C /app/rcon --strip-components=1 \
    && rm /tmp/rcon.tar.gz

# Set environment variables
ENV PORT=3000

# Expose the port
EXPOSE $PORT

# Create a volume for /config
VOLUME /config

# Set the default command to run the binary
CMD ["/app/rcon-cli-web"]
