FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
# Use go work to download dependencies
RUN go work init ./
RUN go work sync

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=0 is important for a static binary, especially in Alpine
# -ldflags="-w -s" strips debug information to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /rbac-worker ./main.go

# --- Final Stage ---
FROM alpine:latest

# It's good practice to run as a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

WORKDIR /app

# Copy the static binary from the builder stage
COPY --from=builder /rbac-worker .

# Set default environment variables.
# LOG_FORMAT can be "console" or "json"
ENV LOG_FORMAT="console"

# Expose a port (optional, as it doesn't run a server, but good practice if it might in the future)
EXPOSE 8080

# Command to run the application
CMD ["./rbac-worker"]