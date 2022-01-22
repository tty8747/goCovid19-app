include .env

all_api: build_api run_api 
all_web: build_web run_web 

init:
	go mod init

get:
	go get -u ./cmd/api

build_api:
	go mod tidy
	mkdir -p ./bin
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_API_NAME}-linux ./cmd/api

build_web:
	go mod tidy
	mkdir -p ./bin
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_WEB_NAME}-linux ./cmd/web

run_api:
	./bin/${BINARY_API_NAME}-linux

run_web:
	./bin/${BINARY_WEB_NAME}-linux

clean_api:
	go clean
	rm ./bin/${BINARY_API_NAME}-linux

clean_web:
	go clean
	rm ./bin/${BINARY_WEB_NAME}-linux
