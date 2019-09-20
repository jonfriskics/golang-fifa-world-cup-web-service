FROM golang:1.13-alpine3.10

ENV CGO_ENABLED 0

WORKDIR /go/src/fifa-world-cup-web-service
COPY . .

RUN go install -v ./...
