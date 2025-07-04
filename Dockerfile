FROM golang:alpine AS builder

RUN apk add --no-cache build-base gcc musl-dev git

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /src

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

ARG VERSION=dev

RUN --mount=type=cache,target=/root/.cache/go-build \
    go build \
      -tags "sqlite_omit_load_extension" \
      -ldflags="-s -w -X main.Version=$VERSION" \
      -o /main .

FROM alpine:latest

RUN apk add --no-cache sqlite-libs ca-certificates \
    && addgroup -S app  \
    && adduser -S app -G app

WORKDIR /app

COPY --from=builder /main /usr/local/bin/main

USER app

CMD ["main", "serve", "--http=0.0.0.0:80"]
