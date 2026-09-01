SERVICES = api ws-server worker outbox-processor topic-creator
IMAGE_TAG ?= dev

.PHONY: run build-go vet test lint check images

run:
	go run ./cmd/api

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
