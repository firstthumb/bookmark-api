FROM golang:1.13.8-alpine3.10

WORKDIR /go/src/app

RUN apk add make git

ADD . .

RUN make build
