include deploy/.env
include deploy/secret.env
LOCAL_BIN:=$(CURDIR)/bin
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

install-deps:
	@if [ ! -f "$(LOCAL_BIN)/protoc-gen-go" ]; then \
		echo "Installing protoc-gen-go..."; \
		GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1; \
	else \
		echo "protoc-gen-go already installed."; \
	fi
	@if [ ! -f "$(LOCAL_BIN)/protoc-gen-go-grpc" ]; then \
		echo "Installing protoc-gen-go-grpc..."; \
		GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2; \
	else \
		echo "protoc-gen-go-grpc already installed."; \
	fi
	@if [ ! -f "$(LOCAL_BIN)/golangci-lint" ]; then \
		echo "Installing golangci-lint..."; \
		GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0; \
	else \
		echo "golangci-lint already installed."; \
	fi
	@if [ ! -f "$(LOCAL_BIN)/goose" ]; then \
    		echo "Installing goose..."; \
    		GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.24.0; \
    	else \
    		echo "goose already installed."; \
    	fi

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate: generate-auth-api

generate-auth-api:
	mkdir -p pkg/chatserver_v1
	protoc --proto_path api/proto/chatserver_v1 \
		--go_out=pkg/chatserver_v1 --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
		--go-grpc_out=pkg/chatserver_v1 --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
		api/proto/chatserver_v1/chatserver.proto

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml


local migration-create:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} create $(name) sql sql

local-migration-status:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v


local-migration-up:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v


local-migration-down:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

build:
	GOOS=linux GOARCH=amd64 go build -o service_linux cmd/grpc_server/main.go
copy-to-server:
	scp service_linux root@$(IP_SERVER):

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t $(REGESTRY)/server-chat:v0.0.1 -f deploy/Dockerfile .
	docker login -u $(USERNAME) -p $(PASSWORD) $(REGESTRY)
	docker push $(REGESTRY)/server-chat:v0.0.1
