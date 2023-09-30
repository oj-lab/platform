OS := $(shell uname -s)

.PHONY: get-front
get-front:
	./script/update-frontend-dist.sh

.PHONY: install-tools
install-tools:
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: build
build: install-tools
	@echo "Building on $(OS)"
	swag fmt -d application/server
	swag init -d application/server -ot go -o application/server/swaggo-gen
	go mod tidy
	go build -o bin/migrate_db application/migrate_db/main.go
	go build -o bin/service application/server/main.go
	go build -o bin/asynq_worker application/asynq_worker/main.go

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

.PHONY: run
run: build check
	make -j run-task-worker run-server

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