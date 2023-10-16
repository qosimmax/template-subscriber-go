.PHONY: build install server test

PROJECT_NAME=$(shell basename $(CURDIR))
PROTO_PATH=$(CURDIR)/proto

## build: build the application
build:
	export GO111MODULE="on"; \
	go mod download; \
	go mod vendor; \
	CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o main cmd/server/main.go

## install: fetches go modules
install:
	export GO111MODULE="on"; \
	go mod tidy; \
	go mod download

## server: runs the server with -race
server:
	export GO111MODULE="on"; \
	go run -race cmd/server/main.go

## test: runs tests
test:
	go test -race ./...

## help: prints help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## docker: Builds and runs the app via the project dockerfile, importing the .env-file as environment variables.
docker:
	docker build -t $(PROJECT_NAME) .
	docker run --rm --name $(PROJECT_NAME) -p 8000:8000 --env-file .env -it $(PROJECT_NAME)


gen-proto:
	protoc \
    --go_out=./example/pb \
    --go_opt=paths=import \
    -I=$(PROTO_PATH)/fake-api \
     $(PROTO_PATH)/fake-api/*.proto \



add-submodule:
	git submodule add git@github.com:qosimmax/services-proto.git ./proto

pull-proto-module:
	git submodule update --recursive --remote