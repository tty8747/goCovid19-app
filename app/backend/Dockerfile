FROM golang:1.17.6-alpine3.15 AS build

WORKDIR /app

ARG CGO_ENABLED=0 \
    GOARCH=amd64 \
    GOOS=linux \
    BINARY_NAME=api \
    PATH_TO_SOURCE=./cmd/api


RUN apk --no-cache add make

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . /app/
RUN make -f Makefile.Docker build

FROM alpine:3.15

ENV DB_HOSTNAME=localhost
ENV MIGRATIONS_DIR=/app/migrations

RUN apk --no-cache add curl

COPY --from=build /app/bin/api-linux /app/api
COPY --from=build /app/migrations /app/migrations
COPY --from=build /app/configs /app/configs

WORKDIR /app

# ARG USERNAME=jarvis \
#     USER_UID=11000 \
#     USER_GID=11000
#
# RUN adduser -DH -u $USER_UID -g $USER_GID $USERNAME \
#     && chown -R $USER_UID:$USER_GID /app
#
# USER $USERNAME:$USER_GID

EXPOSE 5000

CMD ["/app/api"]
