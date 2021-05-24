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

RUN go test ./...
RUN go build -o app ./cmd/archivar/

FROM scratch

COPY --from=builder /build/app /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app", "-c", "/etc/go-archivar/archivar.yaml", "watch"]