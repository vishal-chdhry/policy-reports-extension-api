all: cli server

GOBIN ?= $(shell go env GOPATH)/bin

build-cli:
	cd cli; go build -o ../bin/kubectl-policy-report-extension cmd/main.go

install-cli: build-cli
	cp bin/kubectl-hns-list $(GOBIN)/kubectl-policy-report-extension

cli: install-cli

.PHONY: server
server:
	docker build . -t ghcr.io/vishal-chdhry/policy-reports-aggregation:latest
