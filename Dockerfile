FROM golang:1.11.0-alpine3.8 AS builder
RUN apk add curl git --update --no-cache; \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/tormath1/pub-slack
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only
COPY . ./
RUN GOOS=linux GOARCH=amd64 go build -o pub-slack main.go

FROM alpine:3.8
COPY --from=builder /go/src/github.com/tormath1/pub-slack/pub-slack /usr/local/bin
RUN apk add ca-certificates --update
