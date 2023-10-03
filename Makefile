OS := $(shell uname -s)

.PHONY: get-front
get-front:
	./script/update-frontend-dist.sh

.PHONY: install-tools
install-tools:
	go install github.com/swaggo/swag/cmd/swag@latest
	@# Referencing https://grpc.io/docs/protoc-installation/
	@./script/install-protoc.sh
	@# Track https://grpc.io/docs/languages/go/quickstart/ for update
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

.PHONY: gen-proto
gen-proto: install-tools
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		service/proto/*.proto

.PHONY: build
build: install-tools
	@echo "Building on $(OS)"
	swag fmt -d application/server
	swag init -d application/server -ot go -o application/server/swaggo-gen
	go mod tidy
	go build -o bin/migrate_db application/migrate_db/main.go
	go build -o bin/service application/server/main.go
	go build -o bin/asynq_worker application/asynq_worker/main.go
	go build -o bin/rpc_server application/rpc_server/main.go

.PHONY: clear-db
clear-db:
	docker-compose stop
	docker-compose rm -f

.PHONY: setup-db
setup-db: clear-db build
	docker-compose up -d
	@echo "Wait 10 seconds for db setup"
	sleep 10s
	./bin/migrate_db

.PHONY: check
check:
	go vet ./...

.PHONY: test
test: setup-db check 
	go test -cover -v ./...

.PHONY: run-task-worker
run-task-worker: build check
	./bin/asynq_worker

.PHONY: run-server
run-server: build check
	./bin/service

.PHONY: run-rpc-server
run-rpc-server: build check
	./bin/rpc_server

.PHONY: run
run: build check
	make -j run-task-worker run-server run-rpc-server

.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build     - Build the application"
	@echo "  clear-db  - Clear the database"
	@echo "  setup-db  - Setup the database"
	@echo "  check     - Run go vet"
	@echo "  test      - Run tests"