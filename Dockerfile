FROM --platform=$BUILDPLATFORM golang:1.25.1-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-s -w" \
    -trimpath \
    -o portainer-go-cli \
    .

FROM alpine:latest

RUN apk --no-cache add ca-certificates && \
    adduser -D -s /bin/sh portainer

WORKDIR /home/portainer

COPY --from=builder /build/portainer-go-cli .

RUN chown portainer:portainer portainer-go-cli

USER portainer

ENTRYPOINT ["./portainer-go-cli"]
CMD ["--help"]