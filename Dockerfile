FROM golang:1.23-alpine AS build

RUN apk add --no-cache \
    gcc \
    musl-dev

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

ENV CGO_ENABLED=0
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"
COPY . /workspace/
RUN go build -ldflags='-s -w -extldflags "-static"' -o "ip-syncer"

FROM scratch AS deploy

WORKDIR /app/
COPY --from=build /workspace/ip-syncer /bin/ip-syncer

ENTRYPOINT [ "ip-syncer" ]
