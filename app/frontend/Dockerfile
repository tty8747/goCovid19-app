FROM golang:1.17.6-alpine3.15 AS build

WORKDIR /app

ARG CGO_ENABLED=0 \
    GOARCH=amd64 \
    GOOS=linux \
    BINARY_NAME=web \
    PATH_TO_SOURCE=./cmd/web


RUN apk --no-cache add make

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . /app/
RUN make -f Makefile.Docker build

FROM alpine:3.15

ENV DB_HOSTNAME=localhost
ENV MIGRATIONS_DIR=/app/migrations

COPY --from=build /app/bin/web-linux /app/web
COPY --from=build /app/ui /app/ui

WORKDIR /app

# ARG USERNAME=jarvis \
#     USER_UID=11000 \
#     USER_GID=11000
#
# RUN adduser -DH -u $USER_UID -g $USER_GID $USERNAME \
#     && chown -R $USER_UID:$USER_GID /app
#
# USER $USERNAME:$USER_GID

EXPOSE 4000

CMD ["/app/web", "-addr", "0.0.0.0:4000"]
