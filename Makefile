PORT ?= 80
NAME = kubetimer
REPO_NAME = jijothomasjohn
SVC_VERSION = v0.1

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

.PHONY: publish
publish:
	docker tag $(NAME) $(REPO_NAME)/$(NAME):$(SVC_VERSION)
	docker push $(REPO_NAME)/kubetimer:$(SVC_VERSION)

.PHONY: fix
fix:
	$(HOME)/bin/golangci-lint run --fix

.PHONY: lint
lint:
	$(HOME)/bin/golangci-lint run
