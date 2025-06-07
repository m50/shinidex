# Build stage for JavaScript/TypeScript with Bun
FROM oven/bun:1 AS js-builder

WORKDIR /app

# Copy package files
COPY package.json bun.lockb ./

# Install dependencies
RUN bun install --frozen-lockfile

# Copy frontend source
COPY frontend/ ./frontend/
COPY tailwind.config.js ./
COPY pkg/ ./pkg

# Build the frontend assets
RUN bun run build

# Build stage for Go application
FROM golang:1.24.3 AS go-builder

WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y git gcc && rm -rf /var/lib/apt/lists/*

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download 

# Copy source code
COPY cmd/ ./cmd/
COPY pkg/ ./pkg/
COPY migrations/ ./migrations/

# Copy built frontend assets from js-builder stage
COPY --from=js-builder /app/assets/ ./assets/

# Build the Go application
RUN go generate ./...
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final runtime stage with distroless (includes libgcc)
FROM gcr.io/distroless/cc-debian12

# No need to install anything - distroless comes with ca-certificates

WORKDIR /app/

# Copy the binary from go-builder
COPY --from=go-builder /app/main ./

# Copy static assets and migrations
COPY --from=go-builder /app/assets/ ./assets/
COPY --from=go-builder /app/migrations/ ./migrations/

# Expose port (adjust if your app uses a different port)
EXPOSE 1323

# nonroot
USER nonroot:nonroot

# Run the binary
CMD ["./main"]

