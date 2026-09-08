SERVICES = api ws-server worker outbox-processor topic-creator broadcaster migrate
IMAGE_TAG ?= dev

.PHONY: run run-ws-server build-go vet test lint check images

run:
	set -a && . ./local.env && set +a && go run ./cmd/api

run-ws-server:
	set -a && . ./local.env && set +a && HTTP_PORT=8081 go run ./cmd/ws-server

build-go:
	go build ./...

vet:
	go vet ./...

test:
	go test -race -shuffle=on ./...

lint:
	golangci-lint run ./...

check: build-go vet test lint

images:
	@for s in $(SERVICES); do \
		echo "==> $$s"; \
		docker build --build-arg SERVICE=$$s -t awesome-chat/$$s:$(IMAGE_TAG) . || exit 1; \
	done
