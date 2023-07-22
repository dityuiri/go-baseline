SHELL := /bin/bash

include config/local.env

## Golang Stuff
GOCMD=go
GORUN=$(GOCMD) run
ARGS=$(filter-out $@,$(MAKECMDGOALS))
SRC_PACKAGES=$(shell go list ./...)

ensure-out-dir:
	mkdir -p test_result
deps:
	$(GOCMD) mod tidy

lint: ## lint go code
	golangci-lint run --deadline=30m

fmt: ## format go code
	$(GOCMD) fmt $(SRC_PACKAGES)

run:
	export GOSUMDB=off
	set -o allexport; source config/local.env; set +o allexport && ${GORUN} main.go ${ARGS}

test:
	export GOSUMDB=off
	ginkgo -r

test-coverage: ensure-out-dir
	export GOSUMDB=off
	go test ./... -covermode=count -coverprofile=test_result/coverage-all.out
	go tool cover -html=test_result/coverage-all.out -o test_result/coverage.html
