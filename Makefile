OS := $(shell uname -s)

.PHONY: build
build:
	@echo "Building on $(OS)"
	go mod tidy
	go build -o bin/migrate_user migration/migrate_user.go
	go build -o bin/user_service user-service/application.go

.PHONY: clear-db
clear-db:
	docker-compose stop
	docker-compose rm -f

.PHONY: setup-db
setup-db: build
	docker-compose up -d
	@echo "Wait 10 seconds for db setup"
	sleep 10s
	./bin/migrate_user

.PHONY: check
check:
	go vet ./...

.PHONY: test
test: build check setup-db
	go test -cover -v ./...