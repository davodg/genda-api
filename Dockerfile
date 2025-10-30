FROM golang:1.23.4-bullseye AS builder

ARG BUILD_REF
ARG BUILD_DATE

ENV CGO_ENABLED=0 GO111MODULE=on

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

WORKDIR /go/src/app/cmd/api
RUN go build -ldflags "-s -w -X 'main.buildVersion=${BUILD_REF}'" -o /go/bin/api

FROM debian:bullseye-slim

RUN apt-get update \
 && apt-get install -y --no-install-recommends \
      make \
      ca-certificates \
      curl \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /go/bin/api ./api
COPY Makefile ./
COPY db/migration ./db/migration

EXPOSE 4000

ENTRYPOINT ["./api"]
