FROM golang:alpine AS builder

RUN apk --no-cache add ca-certificates

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -o app ./cmd/archivar/

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to app folder
RUN cp /build/app .

# Build a small image
FROM scratch

COPY --from=builder /dist/app /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Command to run
ENTRYPOINT ["/app"]