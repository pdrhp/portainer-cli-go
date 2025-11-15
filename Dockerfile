FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -o portainer-cli \
    -ldflags="-w -s -X main.version=$(git describe --tags --always --dirty)" \
    ./app

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN adduser -D -s /bin/sh portainer

WORKDIR /home/portainer

COPY --from=builder /app/portainer-cli .

RUN chown portainer:portainer portainer-cli

USER portainer

ENTRYPOINT ["./portainer-cli"]

CMD ["--help"]