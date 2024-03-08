OS := $(shell uname -s)

.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build     - Build the application, swagger document will be generated"
	@echo "  check     - Run go vet"
	@echo "  test      - Run tests, database will be setup"
	@echo "  run       - Run the application"
	@echo "  setup-db  - Setup the database"
	@echo "  clear-db  - Clear the database"
	@echo "  pkg       - Setup the database with problem packages"
	@echo "  all       - Get front dist, Setup-db, read pkg and Run"

.PHONY: build
build: gen-swagger gen-proto
	@echo "Building on $(OS)"
	go mod tidy
	go build -o artifacts/bin/init_db src/application/init_db/main.go
	go build -o artifacts/bin/service src/application/server/main.go
	go build -o artifacts/bin/schedule src/application/schedule/main.go
	go build -o artifacts/bin/read_pkg src/application/read_pkg/main.go

.PHONY: build-image
build-image:
	docker build -f docker/oj-lab-platform.dockerfile -t oj-lab-platform:latest .

.PHONY: clear-db
clear-db:
	docker-compose -f environment/docker-compose.yml -p oj-lab-platform stop
	docker-compose -f environment/docker-compose.yml -p oj-lab-platform rm -f

.PHONY: setup-db
setup-db: clear-db build
	docker-compose -f environment/docker-compose.yml -p oj-lab-platform up -d
	@echo "Wait 10 seconds for db setup"
	sleep 10s
	./artifacts/bin/init_db

.PHONY: gen-swagger
gen-swagger: install-swaggo
	swag fmt -d src/application/server
	swag init -d src/application/server,src/service/model -ot go -o src/application/server/swaggo-gen

.PHONY: check
check: gen-proto
	go vet ./...

.PHONY: test
test: gen-swagger check setup-db
	go test -cover -v -count=1 ./...

.PHONY: pkg
pkg: setup-db
	./artifacts/bin/read_pkg

.PHONY: run-schedule
run-schedule: build check
	./artifacts/bin/schedule

.PHONY: run-server
run-server: build check
	./artifacts/bin/service

.PHONY: run
run: build check
	make -j run-server run-schedule

.PHONY: all
all: get-front pkg
	./artifacts/bin/service

.PHONY: get-front
get-front:
	./scripts/update-frontend-dist.sh artifacts/oj-lab-front/dist

.PHONY: install-swaggo
install-swaggo:
	go install github.com/swaggo/swag/cmd/swag@latest


# Deprecated
# But still needed to pass the build
.PHONY: gen-proto
gen-proto: install-proto
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		src/service/proto/*.proto	

# Deprecated
# But still needed to pass the build
.PHONY: install-proto
install-proto:
	@# Referencing https://grpc.io/docs/protoc-installation/
	@./scripts/install-protoc.sh
	@# Track https://grpc.io/docs/languages/go/quickstart/ for update
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2