# Build stage for frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app

# Copy frontend files
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

# Build stage for backend
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

# Copy go mod files
COPY server/go.mod server/go.sum ./
RUN go mod download

# Copy backend source
COPY server/ ./

# Copy built frontend from previous stage
COPY --from=frontend-builder /app/server/cmd/static ./cmd/static

# Build backend
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gomoku-server ./cmd

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=backend-builder /app/gomoku-server .

# Create data directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080
ENV HOST=0.0.0.0

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# Run the server
CMD ["./gomoku-server"]
