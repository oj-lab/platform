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
	@echo "  gen-swagger 		- Generate swagger document"
	@echo ""
	@echo "Targets for development:"
	@echo "  setup-dependencies  		- Setup the dependencies docker image"
	@echo "  unset-dependencies 		- Unset the dependencies docker image"
	@echo "  get-front 		- Get the frontend files"
	@echo "  check     		- Run go vet"
	@echo "  test      		- Run tests, database will be setup"

.PHONY: build
build: gen-swagger gen-proto
	@echo "Building on $(OS)"
	go mod tidy
	go build -o bin/init_db src/application/init_db/main.go
	go build -o bin/service src/application/server/main.go
	go build -o bin/schedule src/application/schedule/main.go
	go build -o bin/problem_package_loader src/application/problem_package_loader/main.go

.PHONY: gen-swagger
gen-swagger: install-swaggo
	swag fmt -d src/application/server
	swag init -d src/application/server,src/service/model -ot go -o src/application/server/swaggo-gen

# Deprecated
# But still needed to pass the build
.PHONY: gen-proto
gen-proto: install-proto
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		src/service/proto/*.proto	


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
	./bin/problem_package_loader
	docker compose -f $(JUDGER_DOCKER_COMPOSE_FILE) -p oj-lab-judger up -d

.PHONY: get-front
get-front:
	./scripts/update-frontend-dist.sh $(FRONTEND_DIST_DIR)

.PHONY: check
check: gen-proto
	go vet ./...

.PHONY: test
test: gen-swagger check setup-dependencies
	go test -cover -v -count=1 ./...

# Dependent targets

.PHONY: install-swaggo
install-swaggo:
	go install github.com/swaggo/swag/cmd/swag@latest

# Deprecated
# But still needed to pass the build
.PHONY: install-proto
install-proto:
	@# Referencing https://grpc.io/docs/protoc-installation/
	@./scripts/install-protoc.sh
	@# Track https://grpc.io/docs/languages/go/quickstart/ for update
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2