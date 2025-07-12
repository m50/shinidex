FROM oven/bun:1 AS js-builder

WORKDIR /app

COPY package.json bun.lockb ./
RUN bun install --frozen-lockfile

COPY frontend/ ./frontend/
COPY pkg/ ./pkg

RUN bun run build

FROM golang:1.24.3 AS go-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download 

COPY cmd/ ./cmd/
COPY pkg/ ./pkg/

RUN go tool templ generate
RUN go build -o main ./cmd/main/main.go

FROM gcr.io/distroless/cc-debian12

WORKDIR /app

COPY --from=go-builder /app/main ./
COPY --from=js-builder /app/assets/ ./assets/
COPY icons/ ./icons/

EXPOSE 1323
USER nonroot:nonroot
CMD ["./main"]

