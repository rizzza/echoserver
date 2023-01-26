FROM golang:1.19 as builder
ARG IMPORT_PATH=github.com/rizzza/echoserver

COPY . /go/src/$IMPORT_PATH
WORKDIR /go/src/$IMPORT_PATH

ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X github.com/rizzza/echoserver/main.Version=${VERSION}" -a -o echoserver main.go

FROM alpine:latest

COPY --from=builder /go/src/github.com/rizzza/echoserver/echoserver /usr/local/bin
CMD echoserver
