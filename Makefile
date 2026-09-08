SERVICES = api ws-server worker outbox-processor topic-creator broadcaster migrate
IMAGE_TAG ?= dev

LOCAL_BIN = $(CURDIR)/bin
GOWRAP_OP_TPL = $(CURDIR)/pkg/gowrap/tracing.tmpl

GOWRAP_VERSION ?= v1.4.3
MOCKGEN_VERSION ?= v0.6.0
IFACEMAKER_VERSION ?= v1.3.0

.PHONY: run run-ws-server build-go vet test lint check images install-tools gen-mocks gen-interfaces gen-tracing generate

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

install-tools:
	GOBIN=$(LOCAL_BIN) go install github.com/hexdigest/gowrap/cmd/gowrap@$(GOWRAP_VERSION)
	GOBIN=$(LOCAL_BIN) go install go.uber.org/mock/mockgen@$(MOCKGEN_VERSION)
	GOBIN=$(LOCAL_BIN) go install github.com/vburenin/ifacemaker@$(IFACEMAKER_VERSION)

gen-mocks:
	PATH="$(LOCAL_BIN):$(PATH)" go generate -run="mockgen" ./...

gen-interfaces:
	PATH="$(LOCAL_BIN):$(PATH)" go generate -run="ifacemaker" ./...

gen-tracing: gen-interfaces
	PATH="$(LOCAL_BIN):$(PATH)" \
	GOWRAP_OP_TPL="$(GOWRAP_OP_TPL)" \
	go generate -run="gowrap" ./...

generate: gen-mocks gen-tracing
