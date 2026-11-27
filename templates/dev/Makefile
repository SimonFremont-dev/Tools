.PHONY: setup run test lint build clean docker-up docker-down docker-build lint-report

PROJECT ?= {{PROJECT}}
DOCKER_IMAGE ?= $(PROJECT):latest

-include Makefile.lang

SETUP_CMD ?= @echo "setup: nothing configured"
RUN_CMD ?= @echo "run: nothing configured"
TEST_CMD ?= @echo "test: nothing configured"
LINT_CMD ?= @echo "lint: nothing configured"
BUILD_CMD ?= @echo "build: nothing configured"
CLEAN_CMD ?= @echo "clean: nothing configured"

setup:
	$(SETUP_CMD)

run:
	$(RUN_CMD)

test:
	$(TEST_CMD)

lint:
	$(LINT_CMD)

lint-report:
	./scripts/lint-report.sh

build:
	$(BUILD_CMD)

clean:
	$(CLEAN_CMD)

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-build:
	docker build -t $(DOCKER_IMAGE) .
