# Build stage
FROM golang:1.25-trixie AS builder

WORKDIR /app

# Install libvips for building
RUN apt-get update && apt-get install -y --no-install-recommends \
    libvips-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/

# Build the application
RUN CGO_ENABLED=1 go build -o processor ./cmd/processor

# Production stage
FROM debian:trixie-slim AS production

WORKDIR /app

# Install libvips runtime and ca-certificates
RUN apt-get update && apt-get install -y --no-install-recommends \
    libvips \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy binary from builder stage
COPY --from=builder /app/processor .

# Copy config file
COPY configs/processor/config.default.yaml ./configs/

ENTRYPOINT [ "./processor" ]
CMD ["-configs", "configs"]
