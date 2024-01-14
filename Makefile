OS := $(shell uname -s)

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

.PHONY: get-front
get-front:
	./scripts/update-frontend-dist.sh artifacts/oj-lab-front/dist

.PHONY: install-tools
install-tools:
	go install github.com/swaggo/swag/cmd/swag@latest
	@# Referencing https://grpc.io/docs/protoc-installation/
	@./scripts/install-protoc.sh
	@# Track https://grpc.io/docs/languages/go/quickstart/ for update
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

.PHONY: gen-proto
gen-proto: install-tools
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		src/service/proto/*.proto

.PHONY: gen-swagger
gen-swagger: install-tools
	swag fmt -d src/application/server
	swag init -d src/application/server,src/service/model -ot go -o src/application/server/swaggo-gen

.PHONY: build
build: gen-proto gen-swagger
	@echo "Building on $(OS)"
	go mod tidy
	go build -o artifacts/bin/migrate_db src/application/migrate_db/main.go
	go build -o artifacts/bin/service src/application/server/main.go
	go build -o artifacts/bin/asynq_worker src/application/asynq_worker/main.go
	go build -o artifacts/bin/rpc_server src/application/rpc_server/main.go
	go build -o artifacts/bin/schedule src/application/schedule/main.go

.PHONY: clear-db
clear-db:
	docker-compose -f environment/docker-compose.yml -p oj-lab-platform stop
	docker-compose -f environment/docker-compose.yml -p oj-lab-platform rm -f

.PHONY: setup-db
setup-db: clear-db build
	docker-compose -f environment/docker-compose.yml -p oj-lab-platform up -d
	@echo "Wait 10 seconds for db setup"
	sleep 10s
	./artifacts/bin/migrate_db

.PHONY: check
check: gen-proto
	go vet ./...

.PHONY: test
test: gen-swagger check setup-db
	go test -cover -v ./...

.PHONY: run-task-worker
run-task-worker: build check
	./artifacts/bin/asynq_worker

.PHONY: run-schedule
run-schedule: build check
	./artifacts/bin/schedule

.PHONY: run-server
run-server: build check
	./artifacts/bin/service

.PHONY: run-rpc-server
run-rpc-server: build check
	./artifacts/bin/rpc_server

.PHONY: run-background
run-background: build check
	make -j run-schedule

.PHONY: run-all
run-all: build check
	make -j run-server run-schedule

.PHONY: build-docker
build-docker:
	docker build -f docker/oj-lab-platform.dockerfile -t oj-lab-platform:latest .