# Stage 1 - Build the Go application and download CLI
FROM golang:1.19.3-alpine AS builder

# Install necessary build dependencies
RUN apk --no-cache add --update gcc musl-dev

# Set the working directory
WORKDIR /build

# Copy necessary files from the cmd directory
COPY cmd/public /build/public
COPY cmd/routes /build/routes
COPY cmd/config /build/config
COPY cmd/go.mod /build/go.mod
COPY cmd/go.sum /build/go.sum
COPY cmd/main.go /build/main.go

# Build the Go application
ARG VERSION=docker
RUN CGO_ENABLED=1 go build -ldflags "-s -w -X main.ServiceVersion=${VERSION}" -o /build/output/rcon-cli-web /build/*.go

# Stage 2 - Create the final image
FROM alpine AS runner

# Set maintainer label
LABEL maintainer="Xstar97 <dev.xstar97@gmail.com>"

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /build/output/rcon-cli-web /app/

# Create the necessary directories
RUN mkdir -p /config /app/rcon

# Download the latest release of rcon-cli
RUN apk add --no-cache curl tar \
    && mkdir -p /app/rcon \
    && curl -L -o /tmp/rcon.tar.gz $(curl -s https://api.github.com/repos/gorcon/rcon-cli/releases/latest | grep "browser_download_url.*amd64_linux.tar.gz" | cut -d '"' -f 4) \
    && tar -xzf /tmp/rcon.tar.gz -C /app/rcon --strip-components=1 \
    && rm /tmp/rcon.tar.gz

# Create a volume for /config
VOLUME /config

# Set environment variables
ENV PORT=3000

# Expose the port
EXPOSE $PORT

# Set the default command to run the binary
CMD ["/app/rcon-cli-web"]
