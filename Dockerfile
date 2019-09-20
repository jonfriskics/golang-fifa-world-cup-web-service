FROM golang:1.13-alpine3.10

ENV CGO_ENABLED 0

WORKDIR /go/src/golang-fifa-world-cup-web-service
COPY . .

RUN go install -v ./...
