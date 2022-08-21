.PHONY: setup
setup: deps build run

.PHONY: build
build:
	go build cmd/playersInfoGateAwayService.go

.PHONY: run
run:
	go run cmd/playersInfoGateAwayService.go

LOCAL_BIN := $(CURDIR)/bin
export GOBIN := $(LOCAL_BIN)
.PHONY: deps
deps:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc &&