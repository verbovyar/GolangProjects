.PHONY: setup
setup: deps migrationDb generateProto build run

.PHONY: build
build:
	go buid cmd/playersInfoService.go

.PHONY: run
run:
	go run cmd/playersInfoService.go

.PHONY: buf
generateProto:
	buf generate api

export GOOSE_DBSTRING := postgres://postgres:Verbov3232132121@localhost:5432/playersDb?sslmode=disable
export GOOSE_DRIVER := postgres
MIGRATION_DIR=$(CURDIR)/migrations
.PHONY: goose
migrationDb:
	goose -dir=$(MIGRATION_DIR) up

LOCAL_BIN := $(CURDIR)/bin
export GOBIN := $(LOCAL_BIN)
.PHONY: deps
deps:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc && \
    go install github.com/pressly/goose/v3/cmd/goose@latest && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
