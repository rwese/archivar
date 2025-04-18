# Stage 1: Build
FROM golang:1.20-alpine AS builder

# Install necessary dependencies
RUN apk --no-cache add ca-certificates git

# Set environment variables for Go
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Run tests and build the application
RUN go test ./... && \
    go build -o app ./cmd/archivar/

# Stage 2: Final image
FROM scratch

# Copy the built application and certificates
COPY --from=builder /build/app /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Set a non-root user for security
USER 1000:1000

# Set the entrypoint
ENTRYPOINT ["/app", "-c", "/etc/go-archivar/archivar.yaml", "watch"]