SHELL := /bin/bash
# ⚙️ Variables bàsiques
APP_NAME   ?= pinger
IMAGE_REG  ?= ghcr.io/bruli
IMAGE_NAME := $(IMAGE_REG)/$(APP_NAME)
VERSION    ?= 0.1.0
PLATFORM   ?= linux/arm64,linux/amd64
DOCKERFILE ?= Dockerfile
CURRENT_IMAGE := $(IMAGE_NAME):$(VERSION)
CACHE_DIR   ?= .buildx-cache

.PHONY: fmt lint test check clean help security docker-build-image

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

docker-build-image:
	@set -euo pipefail; \
	echo "🐳 Building Docker image $(CURRENT_IMAGE) for ($(PLATFORM))..."; \
	docker buildx build --platform $(PLATFORM) \
		--cache-from type=local,src=$(CACHE_DIR) \
        --cache-to type=local,dest=$(CACHE_DIR),mode=max \
		--build-arg TARGETOS=linux \
		--build-arg TARGETARCH=arm64 \
		-t $(CURRENT_IMAGE) \
		-f $(DOCKERFILE) \
		.

# 🪄 Ajuda
help:
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:' Makefile | awk -F':' '{print "  - " $$1}'
