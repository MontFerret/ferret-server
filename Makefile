.PHONY: build test doc fmt lint vet

export GOPATH

VERSION ?= $(shell git describe --tags --always --dirty)
DIR_BIN = ./bin
DIR_PKG = ./pkg
DIR_SERVER = ./server
DIR_API = ./api
DIR_API_GEN = ${DIR_SERVER}/http/api

default: build

start:
	./bin/ferret-server

build: vet generate test compile

compile:
	go build -v -o ${DIR_BIN}/ferret-server \
	-ldflags "-X main.version=${VERSION}" \
	./main.go

test:
	go test ${DIR_PKG}/...

cover:
	go test -race -coverprofile=coverage.txt -covermode=atomic ${DIR_PKG}/... && \
	curl -s https://codecov.io/bash | bash

generate:
	rm -rf ${DIR_API_GEN} && \
	mkdir ${DIR_API_GEN} && \
	swagger generate server \
		--exclude-main \
		--target=${DIR_API_GEN} \
		--spec=${DIR_API}/api.oas2.json \
		--with-flatten=expand

doc:
	godoc -http=:6060 -index

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt:
	go fmt ${DIR_SERVER}/... ${DIR_PKG}/...

# https://github.com/mgechev/revive
# go get github.com/mgechev/revive
lint:
	revive -config revive.toml \
	-formatter friendly \
	-exclude ./vendor/... \
	-exclude ${DIR_API_GEN}/restapi/... \
	-exclude ${DIR_API_GEN}/models/... \
	./...

# http://godoc.org/code.google.com/p/go.tools/cmd/vet
# go get code.google.com/p/go.tools/cmd/vet
vet:
	go vet ${DIR_SERVER}/... ${DIR_PKG}/...

release:
ifeq ($(RELEASE_VERSION), )
	$(error "Release version is required (version=x)")
else ifeq ($(GITHUB_TOKEN), )
	$(error "GitHub token is required (GITHUB_TOKEN)")
else
	rm -rf ./dist && \
	git tag -a v$(RELEASE_VERSION) -m "New $(RELEASE_VERSION) version" && \
	git push origin v$(RELEASE_VERSION) && \
	goreleaser
endif