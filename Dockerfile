FROM golang:1.13.1-alpine3.10

ENV CGO_ENABLED=0
ENV GOPATH /go
ENV GO111MODULE on

# install git, curl, reflex
RUN apk add --no-cache git curl && \
  curl -fLo ./air https://raw.githubusercontent.com/cosmtrek/air/master/bin/linux/air && chmod +x ./air 