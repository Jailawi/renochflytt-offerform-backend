# --- Build stage ---
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go.mod first to leverage caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY cmd ./cmd
COPY pkg ./pkg
COPY database ./database
COPY templates ./templates
COPY utils ./utils

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

# --- Runtime stage ---
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binary and templates from builder stage
COPY --from=builder /app/app .
COPY --from=builder /app/templates ./templates

# Cloud Run expects port 8080
EXPOSE 8080

ENTRYPOINT ["/app/app"]
