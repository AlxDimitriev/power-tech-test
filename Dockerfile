FROM golang:1.19-alpine3.16 as build-env

WORKDIR /goman/
ADD . .

RUN go mod download
RUN go install ./cmd/...