FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN cd cmd/geoip-api && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /geoip-api .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /geoip-api .

# Create data directory
RUN mkdir -p /data

# Expose port
EXPOSE 8080

# Set default environment variables
ENV GEOIP_DB_PATH=/data/GeoLite2-City.mmdb
ENV PORT=8080
ENV GIN_MODE=release

# Run the application
CMD ["./geoip-api"]
