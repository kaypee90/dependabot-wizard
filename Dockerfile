# Image for the Dependabot Health Report action (see action.yml).
# Also a convenient way to run depbot without a Go toolchain.

FROM golang:1.22-alpine AS build

WORKDIR /src

# Dependencies first, so a source-only change reuses the cached layer.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG VERSION=docker
RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=${VERSION}" -o /out/depbot ./cmd/depbot

FROM alpine:3.20

# Certificates are required to reach api.github.com over TLS.
RUN apk add --no-cache ca-certificates

COPY --from=build /out/depbot /usr/local/bin/depbot

ENTRYPOINT ["/usr/local/bin/depbot"]
