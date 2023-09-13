OS := $(shell uname -s)

.PHONY: build
build:
	@echo "Building on $(OS)"
	go mod tidy
	go build -o bin/migrate_db application/migrate_db/main.go
	go build -o bin/service application/server/main.go

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

.PHONY: run
run: build check
	./bin/service

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