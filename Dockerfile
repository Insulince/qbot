# Use an official lightweight Go image as the base
FROM golang:1.24 as builder

# Set working directory
WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the binary with static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bot -ldflags="-extldflags=-static" ./cmd/qbot

# Use a minimal base image
FROM ubuntu:22.04

# Install dependencies
RUN apt-get update -y && apt-get install -y ca-certificates fuse3 sqlite3

# Copy litefs
COPY --from=flyio/litefs:0.5 /usr/local/bin/litefs /usr/local/bin/litefs

# Create necessary directories for LiteFS
RUN mkdir -p /var/lib/litefs /app

# Set working directory
WORKDIR /app

# Copy the built bot binary from the previous stage
COPY --from=builder /app/bot /app/bot
COPY --from=builder /app/VERSION.md /app/VERSION.md
COPY --from=builder /app/assets /app/assets

# Copy the LiteFS config (you’ll need a litefs.yml file in your repo)
COPY litefs.yml /etc/litefs.yml

# Copy schema and entrypoint script
COPY schema.sql /app/schema.sql
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Set entrypoint
ENTRYPOINT ["/app/entrypoint.sh"]
