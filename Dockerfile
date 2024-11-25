FROM golang:1.23-alpine AS build
ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"

RUN apk add --no-cache \
    gcc \
    musl-dev

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . /workspace/
RUN go build -ldflags='-s -w' -o "ip-syncer"

FROM scratch AS deploy

WORKDIR /app/
COPY --from=build /workspace/ip-syncer /usr/local/bin/ip-syncer

ENTRYPOINT [ "ip-syncer" ]
