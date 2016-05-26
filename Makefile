
GOPACKAGES := $(shell go list ./... | grep -v /vendor/)

BUILD_ENV := CGO_ENABLED=0
BUILD_ARG := -v -x -a -installsuffix cgo
BUILD_DIR := bin


.PHONY: all
all:
	${BUILD_ENV} go build ${BUILD_ARG} -o ${BUILD_DIR}/semver cmd/semver/semver.go

.PHONY: test
test:
	go vet ${GOPACKAGES}
	go test -race -test.v ${GOPACKAGES}

.PHONY: image
image: test all
	docker build --no-cache --tag gcr.io/develo-pe/semver:latest --file Dockerfile .
