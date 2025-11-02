SHELL := /bin/bash

APP_NAME   ?= pinger
IMAGE_REG  ?= ghcr.io/bruli
IMAGE_NAME := $(IMAGE_REG)/$(APP_NAME)
VERSION    ?= 0.1.1
DOCKERFILE ?= Dockerfile
CURRENT_PROD_IMAGE := $(IMAGE_NAME):$(VERSION)
CURRENT_DEV_IMAGE := $(CURRENT_PROD_IMAGE)-dev
OS ?= linux
DEV_ARCH ?= amd64
PROD_ARCH ?= arm64
DEV_PLATFORM := $(OS)/$(DEV_ARCH)
PROD_PLATFORM := $(OS)/$(PROD_ARCH)

.PHONY: fmt lint test check clean help security\
 docker-login  docker-run docker-build-image-dev docker-push-image-prod

.DEFAULT_GOAL := help
# 🧹 Format de codi
fmt:
	@set -euo pipefail; \
	echo "👉 Formating code with gofumpt..."; \
	go tool gofumpt -w .

# 🔍 Linter
lint:
	@set -euo pipefail; \
	echo "🚀 Executing golangci-lint..."; \
	go tool golangci-lint run ./...

# 🧪 Tests amb cobertura i sortida formatejada
test:
	@set -euo pipefail; \
	echo "🚀 Executing tests with cover..."; \
	go test -race ./... -json -cover | go tool tparse -all


# 🧩 Tot en una passada
check: fmt lint security test
	@set -euo pipefail; \
	echo "✅ Format, linter and tests success."

# 🧰 Neteja
clean:
	@set -euo pipefail; \
	echo "🧹 Cleaning cache ..."; \
	go clean -testcache

security:
	@set -euo pipefail; \
	echo "👉 Check security"; \
	go tool govulncheck ./...

docker-login:
	@set -euo pipefail; \
	echo "🔐 Logging into Docker registry..."; \
	echo "$$CR_PAT" | docker login ghcr.io -u bruli --password-stdin

docker-build-image-dev:
	@set -euo pipefail; \
	echo "🐳 Building Docker image $(CURRENT_DEV_IMAGE) for (dev)..."; \
	docker build --platform $(DEV_PLATFORM) \
		--build-arg TARGETOS=$(OS) \
		--build-arg TARGETARCH=$(DEV_ARCH) \
		-t $(CURRENT_DEV_IMAGE) \
		-f $(DOCKERFILE) \
		.
	 echo "✅ Image $(CURRENT_DEV_IMAGE) created successfully."

docker-push-image-prod: docker-login
	@set -euo pipefail; \
	echo "🐳 Building and pushing Docker image $(CURRENT_PROD_IMAGE) for (prod)..."; \
	docker buildx build --platform $(PROD_PLATFORM) \
		--build-arg TARGETOS=$(OS) \
		--build-arg TARGETARCH=$(PROD_ARCH) \
		-t $(CURRENT_PROD_IMAGE) \
		-f $(DOCKERFILE) \
		--load \
		--push \
		.
	 echo "✅ Image $(CURRENT_PROD_IMAGE) pushed successfully."

docker-run: docker-build-image-dev
	@set -euo pipefail; \
    echo "🐳 Running Docker image $(CURRENT_DEV_IMAGE) ..."; \
    docker run --rm -it $(CURRENT_DEV_IMAGE)


# 🪄 Ajuda
help:
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:' Makefile | awk -F':' '{print "  - " $$1}'
