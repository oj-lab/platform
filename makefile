OS := $(shell uname -s)

.PHONY: build
build:
	@echo "Building on $(OS)"
	go mod tidy
	go build -o bin/migrate_user migration/migrate_user.go
	go build -o bin/user_service user-service/application.go

.PHONY: setup-db
setup-db: build
	docker-compose up -d
	 ./bin/migrate_user

.PHONY: test
test: build setup-db
	go test -cover -v ./...