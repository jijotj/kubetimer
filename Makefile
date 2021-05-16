PORT ?= 80
NAME = kubetimer

.PHONY: all
all: lint build

.PHONY: setup
setup:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(HOME)/bin latest

.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	docker build -f ./Dockerfile --build-arg PORT=$(PORT) -t $(NAME) .

.PHONY: fix
fix:
	$(HOME)/bin/golangci-lint run --fix

.PHONY: lint
lint:
	$(HOME)/bin/golangci-lint run
