# Use a minimal builder image with CGO disabled for static linking
FROM golang:1.24 as builder

WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary with static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bot -ldflags="-extldflags=-static" ./cmd/qbot

# Use a lightweight final image (alpine)
FROM alpine:latest

WORKDIR /app

# Copy the statically compiled binary
COPY --from=builder /app/bot /app/bot

# Run the bot
CMD ["/app/bot"]
