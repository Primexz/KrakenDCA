FROM golang:1.22

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . /app/

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /kraken_dca

EXPOSE 8080

# Run
CMD ["/kraken_dca"]