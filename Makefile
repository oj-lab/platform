OS := $(shell uname -s)

FRONTEND_DIST_DIR := frontend/dist
FRONTEND_DIST_URL := https://github.com/oj-lab/oj-lab-front/releases/download/v0.0.3/dist.zip
ICPC_PROBLEM_PACKAGES_DIR := problem_packages/icpc
ICPC_PROBLEM_PACKAGES_URL := https://github.com/oj-lab/problem-packages/releases/download/v0.0.1/icpc_problem.zip

.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build     			- Build the application, swagger document will be generated"
	@echo "  run       			- Run the application"
	@echo "  clean    			- Clean the build"
	@echo "  check     			- Run go vet"
	@echo "  test      			- Run tests, database will be setup"
	@echo "  gen-swagger 			- Generate swagger document"
	@echo "  setup-dependencies  		- Setup the dependencies docker image"
	@echo "  unset-dependencies 		- Unset the dependencies docker image"
	@echo "  get-front 			- Get the frontend files"
	@echo "  update-front 			- Update the frontend files"
	@echo "  get-problem-packages 		- Get the problem packages"
	@echo "  update-problem-packages 	- Update the problem packages"

.PHONY: build
build: gen-swagger gen-proto
	@echo "Building on $(OS)"
	go mod tidy
	go build -o bin/ ./cmd/...

.PHONY: run
run: build
	./bin/web_server

.PHONY: clean
clean:
	rm -rf bin
	rm -rf cmd/web_server/swaggo_gen

.PHONY: gen-swagger
gen-swagger: install-swaggo
	swag fmt -d cmd/web_server
	swag init -d cmd/web_server,models -ot go -o cmd/web_server/swaggo_gen

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
	docker compose stop postgres redis minio clickhouse
	docker compose rm -f postgres redis minio clickhouse

.PHONY: unset-data
unset-data: build
	./bin/clean

.PHONY: setup-dependencies
setup-dependencies: build get-front get-problem-packages
	docker compose up -d postgres redis minio clickhouse
	@echo "Wait 2 seconds for db setup"
	sleep 2s

.PHONY: setup-data
setup-data:setup-dependencies unset-data
	./bin/init

.PHONY: get-front
get-front:
	./scripts/download_and_unzip.sh $(FRONTEND_DIST_DIR) $(FRONTEND_DIST_URL) \
		OVERRIDE=false

.PHONY: update-front
update-front:
	./scripts/download_and_unzip.sh $(FRONTEND_DIST_DIR) $(FRONTEND_DIST_URL) \
		OVERRIDE=true

.PHONY: get-problem-packages
get-problem-packages:
	./scripts/download_and_unzip.sh $(ICPC_PROBLEM_PACKAGES_DIR) $(ICPC_PROBLEM_PACKAGES_URL) \
		OVERRIDE=false

.PHONY: update-problem-packages
update-problem-packages:
	./scripts/download_and_unzip.sh $(ICPC_PROBLEM_PACKAGES_DIR) $(ICPC_PROBLEM_PACKAGES_URL) \
		OVERRIDE=true

.PHONY: check
check: gen-proto install-cilint
	golangci-lint run

.PHONY: test
test: build gen-swagger setup-dependencies
	go test -race -covermode=atomic -coverprofile=coverage.out -cover -v -count=1 \
		./models/... ./modules/... ./services/...

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
	@./scripts/install_protoc.sh
	@# Track https://grpc.io/docs/languages/go/quickstart/ for update
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2