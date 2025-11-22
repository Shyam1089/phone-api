.PHONY: app test build run clean tidy

app:
	@echo "Building and running with Docker..."
	cd deploy && docker-compose up --build

test:
	@echo "Running tests in Docker with coverage..."
	docker build -f deploy/Dockerfile.test -t phone-api-test .
	docker run --rm phone-api-test

build:
	@echo "Building Go binary..."
	go build -o bin/phone-api ./cmd/api

run:
	@echo "Running locally..."
	go run ./cmd/api

clean:
	@echo "Cleaning up..."
	docker rmi -f phone-api-test 2>/dev/null || true
	docker rmi -f deploy-phone-api 2>/dev/null || true
	rm -rf bin/

tidy:
	@echo "Tidying Go modules..."
	go mod tidy