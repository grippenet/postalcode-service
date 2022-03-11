.PHONY: build-data build run docker

INPUT?= laposte_hexasmal.csv
DOCKER_OPTS ?= --rm
VERSION ?= $(shell git describe --tags --abbrev=0)
	
build-data:
	go run cmd/build/build.go $(INPUT)

build:
	go build cmd/server/server.go

run:
	go run cmd/server/server.go

docker:
	docker build -t github.com/grippenet/postalcode-service:$(VERSION)  -f build/docker/Dockerfile $(DOCKER_OPTS) .