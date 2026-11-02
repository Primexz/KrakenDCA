# Step 1: Build the binary
FROM golang:1.23-alpine AS builder

# Install git (needed for commit hash)
RUN apk add --no-cache git

WORKDIR /app
COPY . .

# Download dependencies
RUN go mod download

# Build the binary with commit hash and build date
RUN COMMIT=$(git rev-parse --short HEAD || echo none) && \
    DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ) && \
    go build -ldflags "-X main.commit=$COMMIT -X main.date=$DATE" -o kraken-dca .

# Step 2: Minimal runtime
FROM alpine

COPY --from=builder /app/kraken-dca /usr/bin/kraken-dca

ENTRYPOINT ["/usr/bin/kraken-dca"]
