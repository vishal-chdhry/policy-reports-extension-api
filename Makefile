all: cli server

GOBIN ?= $(shell go env GOPATH)/bin

build-cli:
	cd cli; go build -o ../bin/kubectl-prext cmd/main.go

install-cli: build-cli
	cp bin/kubectl-hns-list $(GOBIN)/kubectl-prext

cli: install-cli

.PHONY: server
server:
	docker build . -t ghcr.io/vishal-chdhry/prext:latest
	docker push ghcr.io/vishal-chdhry/prext:latest
