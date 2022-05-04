FROM golang:1.18.1-alpine3.15 AS goose
# FROM golang:1.17.6-alpine3.15 AS goose

ENV CGO_ENABLED=0
RUN apk update
RUN apk add --no-cache git gcc g++
RUN go install github.com/pressly/goose/cmd/goose@v2.7.0+incompatible
COPY . .
