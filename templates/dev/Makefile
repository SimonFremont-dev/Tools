.PHONY: setup run test lint fmt fmt-check security build clean docker-up docker-down docker-build k8s-deploy k8s-destroy

PROJECT ?= {{PROJECT}}
DOCKER_IMAGE ?= $(PROJECT):latest
KUSTOMIZE_ENV ?= dev

-include Makefile.lang

SETUP_CMD ?= @echo "setup: nothing configured"
RUN_CMD ?= @echo "run: nothing configured"
TEST_CMD ?= @echo "test: nothing configured"
LINT_CMD ?= @echo "lint: nothing configured"
FMT_CMD ?= @echo "fmt: nothing configured"
FMT_CHECK_CMD ?= @echo "fmt-check: nothing configured"
SECURITY_CMD ?= @echo "security: nothing configured"
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

fmt:
	$(FMT_CMD)

fmt-check:
	$(FMT_CHECK_CMD)

security:
	$(SECURITY_CMD)

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

k8s-deploy:
	kubectl apply -k deployments/overlays/$(KUSTOMIZE_ENV)

k8s-destroy:
	kubectl delete -k deployments/overlays/$(KUSTOMIZE_ENV)
