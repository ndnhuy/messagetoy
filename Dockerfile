# module caching
FROM golang:1.21-alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# build
FROM golang:1.21-alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app