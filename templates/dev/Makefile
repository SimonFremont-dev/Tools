.PHONY: setup run test lint build clean docker-up docker-down

PROJECT ?= {{PROJECT}}

setup:
	@echo "Installing dependencies for $(PROJECT)..."

run:
	@echo "Running $(PROJECT)..."

test:
	@echo "Running tests..."

lint:
	@echo "Linting..."

build:
	@echo "Building $(PROJECT)..."

clean:
	@echo "Cleaning..."

docker-up:
	docker compose up -d

docker-down:
	docker compose down
