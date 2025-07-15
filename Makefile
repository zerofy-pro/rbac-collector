GO_PACKAGE := zerofy.pro/rbac-collector
BINARY_NAME := rbac-collector
VERSION ?= $(shell git describe --tags --always --dirty)
APP_VERSION ?= $(subst v,,$(VERSION))

# Docker parameters
DOCKER_REGISTRY ?= ghcr.io
DOCKER_REPO ?= zerofy-pro/rbac-collector
DOCKER_IMAGE := $(DOCKER_REGISTRY)/$(DOCKER_REPO)


# Helm parameters
HELM_CHART_DIR := helm
HELM_CHART_NAME := rbac-collector
HELM_RELEASE_FILE := $(HELM_CHART_NAME)-$(APP_VERSION).tgz

# ====================================================================================
# HELP
# ====================================================================================

.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# ====================================================================================
# DEVELOPMENT
# ====================================================================================

.PHONY: build
build: ## Build the Go binary
	@echo ">> building binary..."
	CGO_ENABLED=0 go build -ldflags="-w -s -X main.version=${VERSION}" -o bin/$(BINARY_NAME) main.go

.PHONY: run
run: ## Run the collector locally using kubeconfig
	@echo ">> running collector locally..."
	LOG_FORMAT=console go run main.go

.PHONY: tidy
tidy: ## Tidy go modules
	@echo ">> tidying go modules..."
	go mod tidy

# ====================================================================================
# DOCKER
# ====================================================================================

.PHONY: docker-build
docker-build: ## Build the Docker image
	@echo ">> building docker image $(DOCKER_IMAGE):$(VERSION)"
	docker build -t $(DOCKER_IMAGE):$(VERSION) .

.PHONY: docker-push-tag
docker-push-tag: ## Push a versioned tag and 'latest'
	@echo ">> pushing docker image $(DOCKER_IMAGE):$(VERSION) and :latest"
	docker push $(DOCKER_IMAGE):$(VERSION)
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest
	docker push $(DOCKER_IMAGE):latest

.PHONY: docker-push-edge
docker-push-edge: ## Push the 'edge' tag for the main branch
	@echo ">> pushing docker image $(DOCKER_IMAGE):edge"
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):edge
	docker push $(DOCKER_IMAGE):edge


# ====================================================================================
# HELM
# ====================================================================================

.PHONY: helm-package
helm-package: ## Package the Helm chart
	@echo ">> packaging helm chart..."
	helm package $(HELM_CHART_DIR) --app-version $(APP_VERSION) --version $(APP_VERSION)

.PHONY: helm-lint
helm-lint: ## Lint the Helm chart
	@echo ">> linting helm chart..."
	helm lint $(HELM_CHART_DIR)

# ====================================================================================
# CI/CD
# ====================================================================================

.PHONY: ci-main-push
ci-main-push: build docker-build docker-push-edge ## CI target for pushes to the main branch

.PHONY: ci-tag-release
ci-tag-release: build docker-build docker-push-tag helm-package ## CI target for creating a tagged release

.PHONY: clean
clean: ## Clean up build artifacts
	@echo ">> cleaning up..."
	@rm -f bin/$(BINARY_NAME)
	@rm -f *.tgz