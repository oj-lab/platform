OS := $(shell uname -s)

DEV_WORKDIR := workdirs/development
DB_DOCKER_COMPOSE_FILE := $(DEV_WORKDIR)/docker-compose.yml
JUDGER_DOCKER_COMPOSE_FILE := $(DEV_WORKDIR)/judger/docker-compose.yml
FRONTEND_DIST_DIR := $(DEV_WORKDIR)/frontend_dist

.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build     		- Build the application, swagger document will be generated"
	@echo "  run       		- Run the application"
	@echo "  clean    		- Clean the build"
	@echo "  check     		- Run go vet"
	@echo "  test      		- Run tests, database will be setup"
	@echo "  gen-swagger 		- Generate swagger document"
	@echo "  setup-dependencies  	- Setup the dependencies docker image"
	@echo "  unset-dependencies 	- Unset the dependencies docker image"
	@echo "  get-front 		- Get the frontend files"

.PHONY: build
build: gen-swagger gen-proto
	@echo "Building on $(OS)"
	go mod tidy
	go build -o bin/init_db cmd/init_db/main.go
	go build -o bin/web_server cmd/web_server/main.go
	go build -o bin/schedule cmd/schedule/main.go
	go build -o bin/problem_loader cmd/problem_loader/main.go

.PHONY: run
run: build
	./bin/web_server

.PHONY: clean
clean:
	rm -rf bin
	rm -rf cmd/web_server/swaggo-gen

.PHONY: gen-swagger
gen-swagger: install-swaggo
	swag fmt -d cmd/web_server
	swag init -d cmd/web_server,models -ot go -o cmd/web_server/swaggo-gen

# Deprecated
# But still needed to pass the build
.PHONY: gen-proto
gen-proto: install-proto
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/*.proto	


# Targets for development

.PHONY: unset-dependencies
unset-dependencies:
	docker compose -f $(JUDGER_DOCKER_COMPOSE_FILE) -p oj-lab-judger stop
	docker compose -f $(JUDGER_DOCKER_COMPOSE_FILE) -p oj-lab-judger rm -f
	docker compose -f $(DB_DOCKER_COMPOSE_FILE) -p oj-lab-dbs stop
	docker compose -f $(DB_DOCKER_COMPOSE_FILE) -p oj-lab-dbs rm -f

.PHONY: setup-dependencies
setup-dependencies: unset-dependencies build
	docker compose -f $(DB_DOCKER_COMPOSE_FILE) -p oj-lab-dbs up -d
	@echo "Wait 10 seconds for db setup"
	sleep 10s
	./bin/init_db
	./bin/problem_loader
	docker compose -f $(JUDGER_DOCKER_COMPOSE_FILE) -p oj-lab-judger up -d

.PHONY: get-front
get-front:
	./scripts/update-frontend-dist.sh $(FRONTEND_DIST_DIR)

.PHONY: check
check: gen-proto install-cilint
	golangci-lint run

.PHONY: test
test: gen-swagger setup-dependencies
	go test -cover -v -count=1 ./...

# Dependent targets

.PHONY: install-swaggo
install-swaggo:
	go install github.com/swaggo/swag/cmd/swag@latest

# See more: https://golangci-lint.run/welcome/install/#local-installation
.PHONY: install-cilint
install-cilint:
	@if [ -z $(shell which golangci-lint) ]; then \
		echo "Installing golangci-lint"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest; \
	fi

# Deprecated
# But still needed to pass the build
.PHONY: install-proto
install-proto:
	@# Referencing https://grpc.io/docs/protoc-installation/
	@./scripts/install-protoc.sh
	@# Track https://grpc.io/docs/languages/go/quickstart/ for update
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2