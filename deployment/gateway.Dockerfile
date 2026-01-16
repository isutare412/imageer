# Build stage
FROM golang:1.25-trixie AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/

# Build the application
RUN CGO_ENABLED=0 go build -o gateway ./cmd/gateway

# Production stage
FROM debian:trixie-slim AS production

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy binary from builder stage
COPY --from=builder /app/gateway .

# Copy config file
COPY configs/gateway/config.default.yaml ./configs/

EXPOSE 8080

ENTRYPOINT [ "./gateway" ]
CMD ["-configs", "configs"]