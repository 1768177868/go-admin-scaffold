APP_NAME=go-admin-base
BUILD_DIR=build
MAIN_FILE=cmd/server/main.go
WORKER_FILE=cmd/worker/main.go

.PHONY: all build clean run test worker

all: build

build:
	@echo "Building server..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "Building worker..."
	@go build -o $(BUILD_DIR)/$(APP_NAME)-worker $(WORKER_FILE)

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

run:
	@echo "Running server..."
	@go run $(MAIN_FILE)

worker:
	@echo "Running worker..."
	@go run $(WORKER_FILE)

test:
	@echo "Running tests..."
	@go test -v ./...

deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

lint:
	@echo "Running linter..."
	@golangci-lint run

dev:
	@echo "Running server in development mode..."
	@air -c .air.toml

.PHONY: migrate
migrate:
	@echo "Running database migrations..."
	@go run cmd/migrate/main.go

.PHONY: seed
seed:
	@echo "Running database seeds..."
	@go run cmd/seed/main.go

.PHONY: docker-build docker-run

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=go-admin-base
MAIN_PATH=cmd/server/main.go

all: lint test build

build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

run:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
	./$(BINARY_NAME)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

lint:
	golangci-lint run

tidy:
	$(GOMOD) tidy

docker-build:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker run -p 8080:8080 $(BINARY_NAME)

migrate-up:
	go run cmd/migration/main.go up

migrate-down:
	go run cmd/migration/main.go down

dev:
	air -c .air.toml

help:
	@echo "make - Build and run tests"
	@echo "make build - Build the binary"
	@echo "make run - Run the application"
	@echo "make test - Run tests"
	@echo "make clean - Clean build files"
	@echo "make lint - Run linter"
	@echo "make tidy - Tidy go.mod"
	@echo "make docker-build - Build Docker image"
	@echo "make docker-run - Run Docker container"
	@echo "make migrate-up - Run database migrations"
	@echo "make migrate-down - Rollback database migrations"
	@echo "make dev - Run with hot reload" 