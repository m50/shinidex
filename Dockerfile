# Build stage for JavaScript/TypeScript with Bun
FROM oven/bun:1 AS js-builder

WORKDIR /app

COPY package.json bun.lockb ./
RUN bun install --frozen-lockfile

COPY frontend/ ./frontend/
COPY tailwind.config.js ./
COPY pkg/ ./pkg

RUN bun run build

FROM golang:1.24.3 AS go-builder

WORKDIR /app

RUN apt-get update \
    && apt-get install -y git gcc \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download 

COPY cmd/ ./cmd/
COPY pkg/ ./pkg/
COPY migrations/ ./migrations/

RUN go generate ./...
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

FROM gcr.io/distroless/cc-debian12

WORKDIR /app

COPY --from=go-builder /app/main ./
COPY --from=js-builder /app/assets/ ./assets/
COPY --from=go-builder /app/migrations/ ./migrations/
COPY icons/ ./icons/

EXPOSE 1323
USER nonroot:nonroot
CMD ["./main"]

